package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/slug"
)

var (
	ErrProductNameAlreadyTaken = errors.New("given product name is already in use")
	ErrProductNotFound         = errors.New("product not found")
)

type productService struct {
	log         *slog.Logger
	productRepo contract.ProductRepository
}

func NewProductService(log *slog.Logger, productRepo contract.ProductRepository) contract.ProductService {
	return &productService{log, productRepo}
}

func (svc *productService) UpdatePrice(publicID string, price float64) error {
	storedProduct, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error(fmt.Sprintf("product with publicID [%s] not found", publicID))
		return err
	}

	p, err := product.New(storedProduct.Name, storedProduct.Price, storedProduct.Quantity)
	if err != nil {
		svc.log.Error(err.Error())
		return err
	}

	if err = p.ChangePrice(price); err != nil {
		svc.log.Error(fmt.Sprintf("error changing product price %s", err.Error()))
		return err
	}

	err = svc.productRepo.Update(product.Product{PublicID: publicID, Price: p.Price})
	if err != nil {
		svc.log.Error(err.Error())
		return errors.New("cannot increment product price")
	}
	return nil
}

func (svc *productService) Create(data dto.CreateProduct) error {
	p, err := product.New(data.Name, data.Price, data.Quantity)
	if err != nil {
		svc.log.Error(err.Error())
		return err
	}

	_, err = svc.productRepo.Create(*p)
	if err != nil {
		if strings.Contains(err.Error(), "uni_products_slug") {
			svc.log.Error(fmt.Sprintf("the product name [%s] is already in use", p.Name))
			return ErrProductNameAlreadyTaken
		}
		svc.log.Error(err.Error())
		return err
	}
	return nil
}

func (svc *productService) Delete(publicID string) error {
	_, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error(fmt.Sprintf("product with PublicID [%s] not found", publicID))
		return ErrProductNotFound
	}

	if err = svc.productRepo.Delete(publicID); err != nil {
		return err
	}
	return nil
}

func (svc *productService) GetAll() ([]product.Product, error) {
	products, err := svc.productRepo.GetAll()
	if err != nil {
		svc.log.Error("cannot get products slice")
		return nil, err
	}
	return products, nil
}

func (svc *productService) GetByPublicID(publicID string) (*product.Product, error) {
	p, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error(fmt.Sprintf("product with PublicID [%s] not found", publicID))
		return nil, err
	}
	return p, nil
}

func (svc *productService) GetBySlug(slugInput string) (*product.Product, error) {
	s, err := slug.New(slugInput)
	if err != nil {
		svc.log.Error(fmt.Sprintf("invalid slug [%s]", string(s.GetSlug())))
		return nil, err
	}

	p, err := svc.productRepo.GetBySlug(string(s.GetSlug()))
	if err != nil {
		svc.log.Error(fmt.Sprintf("product with slug [%s] not found", string(s.GetSlug())))
		return nil, err
	}
	return p, nil
}
