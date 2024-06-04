package service

import (
	"log/slog"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"
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
		return err
	}

	_, err = p.productRepo.Create(*product)
	if err != nil {
		return err
	}

	return nil
}
