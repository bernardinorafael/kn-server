package gormrepo

import (
	"fmt"
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"gorm.io/gorm"
)

/*
* TODO: remove entity/model mapping logic from repositories and do it into service layer
 */

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) contract.ProductRepository {
	return &productRepo{db}
}

func (r productRepo) Update(p product.Product) (gormodel.Product, error) {
	var product gormodel.Product

	err := r.db.
		Where("public_id = ?", p.PublicID()).
		First(&product).
		Error
	if err != nil {
		return product, err
	}

	product.Name = p.Name()
	product.Slug = string(p.Slug())
	product.Price = int(p.Price())
	product.Quantity = p.Quantity()
	product.Enabled = p.Enabled()
	product.UpdatedAt = time.Now()

	err = r.db.Save(&product).Error
	if err != nil {
		return gormodel.Product{}, err
	}

	return product, nil
}

func (r productRepo) Create(p product.Product) (gormodel.Product, error) {
	var product gormodel.Product

	newProduct := gormodel.Product{
		Name:      p.Name(),
		Price:     int(p.Price()),
		Quantity:  p.Quantity(),
		PublicID:  p.PublicID(),
		Slug:      string(p.Slug()),
		Image:     p.Image(),
		Enabled:   p.Enabled(),
		UpdatedAt: time.Now(),
	}

	err := r.db.
		Create(newProduct).
		First(&product).
		Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r productRepo) GetByID(id int) (gormodel.Product, error) {
	var product gormodel.Product

	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r productRepo) GetBySlug(slug string) (gormodel.Product, error) {
	var product gormodel.Product

	err := r.db.Where("slug = ?", slug).First(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r productRepo) Delete(publicID string) error {
	var product gormodel.Product

	err := r.db.
		Where("public_id = ?", publicID).
		First(&product).
		Delete(&product).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r productRepo) GetAll(disabled bool, orderBy string) ([]gormodel.Product, error) {
	var products []gormodel.Product
	var enabled bool = !disabled

	err := r.db.
		Where("enabled = ?", enabled).
		Order(fmt.Sprintf("%v desc", orderBy)).
		Find(&products).Error
	if err != nil {
		return products, err
	}

	return products, nil
}

func (r productRepo) GetByPublicID(publicID string) (gormodel.Product, error) {
	var product gormodel.Product

	err := r.db.
		Where("public_id = ?", publicID).
		First(&product).
		Error
	if err != nil {
		return product, err
	}

	return product, nil
}
