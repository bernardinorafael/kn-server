package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	env "github.com/bernardinorafael/kn-server/internal/config"
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/slug"
	"github.com/bernardinorafael/kn-server/internal/infra/s3client"
)

var (
	ErrProductNameAlreadyTaken = errors.New("given product name is already in use")
	ErrProductNotFound         = errors.New("product not found")
)

type productService struct {
	log         *slog.Logger
	env         *env.Env
	productRepo contract.ProductRepository
}

func NewProductService(log *slog.Logger, env *env.Env, productRepo contract.ProductRepository) contract.ProductService {
	return &productService{log, env, productRepo}
}

func (svc *productService) IncreaseQuantity(publicID string, quantity int32) error {
	storedProduct, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error(fmt.Sprintf("product with publicID %s not found", publicID))
		return err
	}

	p, err := product.New(storedProduct.Name, storedProduct.Price, storedProduct.Quantity)
	if err != nil {
		svc.log.Error(err.Error())
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
	storedProduct, err := svc.productRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error(fmt.Sprintf("product with publicID %s not found", publicID))
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

func (svc *productService) Create(data dto.CreateProduct, file io.Reader, fileName string) error {
	p, err := product.New(data.Name, data.Price, data.Quantity)
	if err != nil {
		svc.log.Error(err.Error())
		return err
	}

	ext := filepath.Ext(fileName)
	if len(ext) == 0 {
		err = errors.New("file has no extension")
		svc.log.Error(err.Error())
		return err
	}

	client, err := s3client.New(svc.env)
	if err != nil {
		svc.log.Error(fmt.Sprintf("cannot initialize aws s3 service %v", err))
		return errors.New("cannot initialize aws s3 service")
	}

	uploader := manager.NewUploader(client)
	res, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(svc.env.AWSBucket),
		Key:         aws.String(fmt.Sprintf("%s%s", p.PublicID, ext)),
		ContentType: aws.String("image/*"),
		ACL:         "public-read",
		Body:        file,
	})

	p.SetImageURL(res.Location)

	if err != nil {
		svc.log.Error(fmt.Sprintf("cannot initialize aws s3 service %v", err))
		return errors.New("cannot init s3")
	}

	_, err = svc.productRepo.Create(*p)
	if err != nil {
		if strings.Contains(err.Error(), "uni_products_slug") {
			svc.log.Error(fmt.Sprintf("the product name %s is already in use", p.Name))
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
		svc.log.Error(fmt.Sprintf("product with PublicID %s not found", publicID))
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
		svc.log.Error(fmt.Sprintf("product with PublicID %s not found", publicID))
		return nil, err
	}
	return p, nil
}

func (svc *productService) GetBySlug(slugInput string) (*product.Product, error) {
	s, err := slug.New(slugInput)
	if err != nil {
		svc.log.Error(fmt.Sprintf("invalid slug %s", string(s.GetSlug())))
		return nil, err
	}

	p, err := svc.productRepo.GetBySlug(string(s.GetSlug()))
	if err != nil {
		svc.log.Error(fmt.Sprintf("product with slug %s not found", string(s.GetSlug())))
		return nil, err
	}
	return p, nil
}
