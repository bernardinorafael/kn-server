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
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/google/uuid"
)

var (
	ErrProductNameAlreadyTaken = errors.New("given product name is already in use")
	ErrProductNotFound         = errors.New("product not found")
)

type WithProductParams struct {
	Log         logger.Logger
	Env         *env.Env
	ProductRepo contract.ProductRepository
	FileService contract.FileManagerService
}

type productService struct {
	log         logger.Logger
	env         *env.Env
	productRepo contract.ProductRepository
	fileService contract.FileManagerService
}

func NewProductService(p WithProductParams) contract.ProductService {
	return &productService{
		log:         p.Log,
		env:         p.Env,
		productRepo: p.ProductRepo,
		fileService: p.FileService,
	}
}

func (svc productService) ChangeStatus(publicID string, status bool) error {
	found, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("product not found", "public_id", publicID)
		return err
	}

	p, err := product.New(product.Params{
		PublicID: found.PublicID,
		Name:     found.Name,
		Image:    found.Image,
		Price:    found.Price,
		Quantity: found.Quantity,
		Enabled:  found.Enabled,
	})
	if err != nil {
		svc.log.Error("failed to initialize new product entity", "error", err.Error())
		return err
	}
	p.ChangeStatus(status)

	productModel := gormodel.Product{
		PublicID: p.PublicID(),
		Slug:     string(p.Slug()),
		Name:     p.Name(),
		Image:    p.Image(),
		Price:    int(p.Price()),
		Quantity: p.Quantity(),
		Enabled:  p.Enabled(),
	}

	_, err = svc.productRepo.Update(productModel)
	if err != nil {
		svc.log.Error("cannot change product status", "error", err.Error())
		return errors.New("cannot change product status")
	}

	return nil
}

func (svc productService) IncreaseQuantity(publicID string, quantity int) error {
	found, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("product not found", "public_id", publicID)
		return err
	}

	p, err := product.New(product.Params{
		PublicID: found.PublicID,
		Name:     found.Name,
		Image:    found.Image,
		Price:    found.Price,
		Quantity: found.Quantity,
		Enabled:  found.Enabled,
	})
	if err != nil {
		svc.log.Error("failed to initialize new product entity", "error", err.Error())
		return err
	}

	if err = p.IncreaseQuantity(quantity); err != nil {
		svc.log.Error(fmt.Sprintf("failed increment product quantity %s", err.Error()))
		return err
	}

	productModel := gormodel.Product{
		PublicID: p.PublicID(),
		Slug:     string(p.Slug()),
		Name:     p.Name(),
		Image:    p.Image(),
		Price:    int(p.Price()),
		Quantity: p.Quantity(),
		Enabled:  p.Enabled(),
	}

	_, err = svc.productRepo.Update(productModel)
	if err != nil {
		svc.log.Error(err.Error())
		return errors.New("cannot increment product price")
	}

	return nil
}

func (svc productService) UpdatePrice(publicID string, price int) error {
	found, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("product not found", "public_id", publicID)
		return err
	}

	p, err := product.New(product.Params{
		PublicID: found.PublicID,
		Name:     found.Name,
		Image:    found.Image,
		Price:    found.Price,
		Quantity: found.Quantity,
		Enabled:  found.Enabled,
	})
	if err != nil {
		svc.log.Error("failed to initialize new product entity", "error", err.Error())
		return err
	}

	if err = p.ChangePrice(price); err != nil {
		svc.log.Error("cannot change product price", "error", err.Error())
		return err
	}

	productModel := gormodel.Product{
		PublicID: p.PublicID(),
		Slug:     string(p.Slug()),
		Name:     p.Name(),
		Image:    p.Image(),
		Price:    int(p.Price()),
		Quantity: p.Quantity(),
		Enabled:  p.Enabled(),
	}

	_, err = svc.productRepo.Update(productModel)
	if err != nil {
		svc.log.Error("error updating product", "error", err.Error())
		return errors.New("cannot increment product price")
	}
	return nil
}

func (svc productService) Create(dto dto.CreateProduct) error {
	id := uuid.NewString()

	ext := filepath.Ext(dto.ImageName)
	if len(ext) == 0 {
		svc.log.Error("image name cannot be empty")
		return errors.New("cannot reach image name")
	}
	filename := fmt.Sprintf("%s%s", id, ext)

	location, err := svc.fileService.UploadFile(dto.Image, filename, svc.env.AWSBucket)
	if err != nil {
		svc.log.Error("cannot upload image to s3", "error", err.Error())
		return err
	}

	p, err := product.New(product.Params{
		PublicID: id,
		Name:     dto.Name,
		Image:    location,
		Price:    dto.Price,
		Quantity: dto.Quantity,
		Enabled:  true,
	})
	if err != nil {
		svc.log.Error("failed to initialize new product entity", "error", err.Error())
		return err
	}

	productModel := gormodel.Product{
		PublicID: p.PublicID(),
		Slug:     string(p.Slug()),
		Name:     p.Name(),
		Image:    p.Image(),
		Price:    int(p.Price()),
		Quantity: p.Quantity(),
		Enabled:  p.Enabled(),
	}

	err = svc.productRepo.Create(productModel)
	if err != nil {
		if strings.Contains(err.Error(), "uni_products_slug") {
			svc.log.Error("product name already taken", "name", dto.Name)
			return ErrProductNameAlreadyTaken
		}
		svc.log.Error(err.Error())
		return err
	}

	return nil
}

func (svc productService) Delete(publicID string) error {
	_, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("product not found", "public_id", publicID)
		return ErrProductNotFound
	}

	if err = svc.productRepo.Delete(publicID); err != nil {
		svc.log.Error("error deleting product", "error", err.Error())
		return errors.New("error deleting product")
	}

	return nil
}

func (svc productService) GetAll(disabled bool, orderBy string) ([]gormodel.Product, error) {
	products, err := svc.productRepo.GetAll(disabled, orderBy)

	if err != nil {
		svc.log.Error("error retrieving products", "errors", err.Error())
		return nil, errors.New("products not found")
	}

	return products, nil
}

func (svc productService) GetByPublicID(publicID string) (gormodel.Product, error) {
	var product gormodel.Product

	p, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("product not found", "public_id", publicID)
		return product, errors.New("product not found")
	}

	return p, nil
}

func (svc productService) GetBySlug(slugInput string) (gormodel.Product, error) {
	var product gormodel.Product

	s, err := slug.New(slugInput)
	if err != nil {
		svc.log.Error("slug validation error", "error", err.Error())
		return product, err
	}

	p, err := svc.productRepo.GetBySlug(string(s.GetSlug()))
	if err != nil {
		svc.log.Error("product not found", "slug", slugInput)
		return product, errors.New("product not found")
	}

	return p, nil
}
