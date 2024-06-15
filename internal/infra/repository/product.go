package repository

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"

	"gorm.io/gorm"
)

type productRepo struct {
	DB *gorm.DB
}

func NewProductRepo(DB *gorm.DB) contract.ProductRepository {
	return &productRepo{DB}
}

func (p *productRepo) Create(prod product.Product) (*product.Product, error) {
	newProduct := product.Product{
		Name:     prod.Name,
		Price:    prod.Price,
		Quantity: prod.Quantity,
	}

	err := p.DB.Create(&prod).Error
	if err != nil {
		return nil, err
	}
	return &newProduct, nil
}

func (p *productRepo) FindByID(id int) (*product.Product, error) {
	var product product.Product

	err := p.DB.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *productRepo) FindBySlug(name string) (*product.Product, error) {
	var product product.Product

	err := p.DB.Where("slug = ?", name).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *productRepo) Delete(id int) error {
	err := p.DB.Delete(&product.Product{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
