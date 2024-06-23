package repository

import (
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) contract.ProductRepository {
	return &productRepo{db}
}

func (p *productRepo) Create(prod product.Product) (*product.Product, error) {
	newProduct := product.Product{
		Name:     prod.Name,
		Price:    prod.Price,
		Quantity: prod.Quantity,
	}

	err := p.db.Create(&prod).Error
	if err != nil {
		return nil, err
	}
	return &newProduct, nil
}

func (p *productRepo) GetByID(id int) (*product.Product, error) {
	var product product.Product

	err := p.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *productRepo) GetBySlug(name string) (*product.Product, error) {
	var product product.Product

	err := p.db.Where("slug = ?", name).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *productRepo) Delete(publicID string) error {
	var product product.Product

	err := p.db.Where("public_id = ?", publicID).Delete(&product).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *productRepo) GetAll() ([]product.Product, error) {
	products := make([]product.Product, 0)

	err := p.db.Find(&products).Error
	if err != nil {
		return products, err
	}

	return products, nil
}

func (p *productRepo) GetByPublicID(publicID string) (*product.Product, error) {
	var product product.Product

	err := p.db.Where("public_id = ?", publicID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
