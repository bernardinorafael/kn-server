package repository

import (
	"errors"
	"time"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/gorm"
)

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) contract.ProductRepository {
	return &productRepository{DB}
}

func (pr *productRepository) GetByName(p entity.Product) (*entity.Product, error) {
	product := entity.Product{}

	err := pr.DB.Where("name = ?", p.Name).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (pr productRepository) Save(p entity.Product) error {
	product := entity.Product{
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		SKU:         p.SKU,
		Size:        p.Size,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := pr.DB.Create(&product).Error; err != nil {
		return err
	}

	return nil
}

func (pr productRepository) GetAll() (*[]entity.Product, error) {
	return nil, nil
}
