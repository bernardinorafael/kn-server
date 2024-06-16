package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/slug"
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

func (p *productService) Create(data dto.CreateProduct) error {
	product, err := product.New(data.Name, data.Price, data.Quantity)
	if err != nil {
		p.log.Error(err.Error())
		return err
	}

	_, err = p.productRepo.Create(*product)
	if err != nil {
		if strings.Contains(err.Error(), "uni_products_slug") {
			p.log.Error(fmt.Sprintf("the product name [%s] is already in use", product.Name))
			return ErrProductNameAlreadyTaken
		}
		p.log.Error(err.Error())
		return err
	}
	return nil
}

func (p *productService) Delete(publicID string) error {
	_, err := p.productRepo.GetByPublicID(publicID)
	if err != nil {
		p.log.Error(fmt.Sprintf("product with PublicID [%s] not found", publicID))
		return ErrProductNotFound
	}

	if err = p.productRepo.Delete(publicID); err != nil {
		return err
	}
	return nil
}

func (p *productService) GetAll() ([]product.Product, error) {
	products := make([]product.Product, 0)

	allProducts, err := p.productRepo.GetAll()
	if err != nil {
		p.log.Error("cannot get products slice")
		return nil, err
	}
	products = allProducts

	return products, nil
}

func (p *productService) GetByPublicID(publicID string) (*product.Product, error) {
	product, err := p.productRepo.GetByPublicID(publicID)
	if err != nil {
		p.log.Error(fmt.Sprintf("product with PublicID [%s] not found", publicID))
		return nil, err
	}
	return product, nil
}

func (p *productService) GetBySlug(slugInput string) (*product.Product, error) {
	s, err := slug.New(slugInput)
	if err != nil {
		p.log.Error(fmt.Sprintf("invalid slug [%s]", string(s.GetSlug())))
		return nil, err
	}

	product, err := p.productRepo.GetBySlug(string(s.GetSlug()))
	if err != nil {
		p.log.Error(fmt.Sprintf("product with slug [%s] not found", string(s.GetSlug())))
		return nil, err
	}
	return product, nil
}
