package service

import (
	"log/slog"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type productService struct {
	l           *slog.Logger
	productRepo contract.ProductRepository
}

func NewProductService(l *slog.Logger, productRepo contract.ProductRepository) contract.ProductService {
	return &productService{l, productRepo}
}

func (p *productService) Create(data dto.CreateProduct) error {
	product, err := entity.NewProduct(data.Name, data.Price, data.Quantity)
	if err != nil {
		return err
	}

	_, err = p.productRepo.Create(*product)
	if err != nil {
		return err
	}

	return nil
}
