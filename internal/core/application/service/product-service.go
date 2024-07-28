package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	env "github.com/bernardinorafael/kn-server/internal/config"
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/slug"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
	"github.com/bernardinorafael/kn-server/pkg/logger"
)

var (
	ErrProductNameAlreadyTaken = errors.New("given product name is already in use")
	ErrProductNotFound         = errors.New("product not found")
)

type productService struct {
	log         logger.Logger
	env         *env.Env
	productRepo contract.ProductRepository
	fileService contract.FileManagerService
}

func NewProductService(log logger.Logger, env *env.Env, productRepo contract.ProductRepository, fileService contract.FileManagerService) contract.ProductService {
	return &productService{log, env, productRepo, fileService}
}

func (svc *productService) ChangeStatus(publicID string, status bool) error {
	record, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("product not found", "public_id", publicID)
		return err
	}

	p, err := product.New(record.Name, record.Price, record.Quantity)
	if err != nil {
		svc.log.Error("cannot init product entity", "error", err.Error())
		return err
	}
	p.ChangeStatus(status)

	updated := product.Product{
		PublicID: record.PublicID,
		Name:     p.Name,
		Slug:     p.Slug,
		Price:    p.Price,
		Quantity: p.Quantity,
		Enabled:  p.GetStatus(),
	}

	err = svc.productRepo.Update(updated)
	if err != nil {
		svc.log.Error("cannot change product status", "error", err.Error())
		return errors.New("cannot change product status")
	}
	return nil
}

func (svc *productService) IncreaseQuantity(publicID string, quantity int32) error {
	record, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("product not found", "public_id", publicID)
		return err
	}

	p, err := product.New(record.Name, record.Price, record.Quantity)
	if err != nil {
		svc.log.Error("cannot init product entity", "error", err.Error())
		return err
	}

	if err = p.IncreaseQuantity(quantity); err != nil {
		svc.log.Error(fmt.Sprintf("error increment product quantity %s", err.Error()))
		return err
	}

	err = svc.productRepo.Update(product.Product{PublicID: publicID, Quantity: p.Quantity})
	if err != nil {
		svc.log.Error(err.Error())
		return errors.New("cannot increment product price")
	}
	return nil
}

func (svc *productService) UpdatePrice(publicID string, price float64) error {
	record, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("product not found", "public_id", publicID)
		return err
	}

	p, err := product.New(record.Name, record.Price, record.Quantity)
	if err != nil {
		svc.log.Error(err.Error())
		return err
	}

	if err = p.ChangePrice(price); err != nil {
		svc.log.Error("cannot change product price", "error", err.Error())
		return err
	}

	err = svc.productRepo.Update(product.Product{PublicID: publicID, Price: p.Price})
	if err != nil {
		svc.log.Error("error updating product", "error", err.Error())
		return errors.New("cannot increment product price")
	}
	return nil
}

func (svc *productService) Create(data dto.CreateProduct) error {
	product, err := product.New(data.Name, data.Price, data.Quantity)
	if err != nil {
		svc.log.Error("invalid product entity", "error", err.Error())
		return err
	}

	ext := filepath.Ext(data.ImageName)
	if len(ext) == 0 {
		svc.log.Error("image name cannot be empty")
		return errors.New("cannot reach image name")
	}
	filename := fmt.Sprintf("%s%s", product.PublicID, ext)

	res, err := svc.fileService.UploadFile(data.Image, filename, svc.env.AWSBucket)
	if err != nil {
		svc.log.Error("cannot upload image to s3", "error", err.Error())
		return err
	}
	product.SetImage(res.Location)

	_, err = svc.productRepo.Create(*product)
	if err != nil {
		if strings.Contains(err.Error(), "uni_products_slug") {
			svc.log.Error("product name already taken", "name", data.Name)
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
		svc.log.Error("product not found", "public_id", publicID)
		return ErrProductNotFound
	}
	if err = svc.productRepo.Delete(publicID); err != nil {
		return err
	}
	return nil
}

func (svc *productService) GetAll(disabled bool, orderBy string) ([]model.Product, error) {
	products, err := svc.productRepo.GetAll(disabled, orderBy)
	if err != nil {
		svc.log.Error("cannot get products slice")
		return nil, err
	}
	return products, nil
}

func (svc *productService) GetByPublicID(publicID string) (*model.Product, error) {
	p, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("product not found", "public_id", publicID)
		return nil, err
	}
	return p, nil
}

func (svc *productService) GetBySlug(slugInput string) (*model.Product, error) {
	s, err := slug.New(slugInput)
	if err != nil {
		svc.log.Error("cannot init slug value object", "error", err.Error())
		return nil, err
	}

	p, err := svc.productRepo.GetBySlug(string(s.GetSlug()))
	if err != nil {
		svc.log.Error("not found a product by given slug", "slug", slugInput)
		return nil, err
	}
	return p, nil
}
