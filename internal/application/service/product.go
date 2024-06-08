package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"
)

var (
	ErrProductNameAlreadyTaken = errors.New("given product name is already in use")
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
