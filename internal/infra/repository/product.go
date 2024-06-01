package repository

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/gorm"
)

type productRepo struct {
	DB *gorm.DB
}

func NewProductRepo(DB *gorm.DB) contract.ProductRepository {
	return &productRepo{DB}
}

func (p *productRepo) Create(product entity.Product) (*entity.Product, error) {
	prod := &entity.Product{
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}

	err := p.DB.Create(&prod).Error
	if err != nil {
		return nil, err
	}

	return prod, nil
}
