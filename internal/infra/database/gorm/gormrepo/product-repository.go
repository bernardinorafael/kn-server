package gormrepo

import (
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) contract.ProductRepository {
	return &productRepo{db}
}

func (r *productRepo) Update(prod product.Product) error {
	var p product.Product

	updated := model.Product{
		Name:      prod.Name,
		Slug:      prod.Slug,
		Price:     prod.Price,
		Quantity:  prod.Quantity,
		Enabled:   prod.Enabled,
		UpdatedAt: time.Now(),
	}

	err := r.db.
		Model(&p).
		Where("public_id = ?", prod.PublicID).
		Updates(updated).
		Error

	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) Create(prod product.Product) (*model.Product, error) {
	newProduct := &model.Product{
		Name:     prod.Name,
		Price:    prod.Price,
		Quantity: prod.Quantity,
	}

	err := r.db.Create(&prod).Error
	if err != nil {
		return nil, err
	}
	return newProduct, nil
}

func (r *productRepo) GetByID(id int) (*model.Product, error) {
	var product model.Product

	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) GetBySlug(name string) (*model.Product, error) {
	var product model.Product

	err := r.db.Where("slug = ?", name).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) Delete(publicID string) error {
	var product product.Product

	err := r.db.Where("public_id = ?", publicID).Delete(&product).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) GetAll() ([]model.Product, error) {
	products := make([]model.Product, 0)

	err := r.db.Order("created_at desc").Find(&products).Error
	if err != nil {
		return products, err
	}
	return products, nil
}

func (r *productRepo) GetByPublicID(publicID string) (*model.Product, error) {
	var product model.Product

	err := r.db.Where("public_id = ?", publicID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
