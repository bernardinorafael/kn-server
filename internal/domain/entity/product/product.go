package product

import (
	"errors"
	"fmt"
	"time"

	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/slug"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	minNameLength = 3
	maxNameLength = 120
)

var (
	ErrInvalidQuantity    = errors.New("product quantity cannot be zero")
	ErrInvalidPrice       = errors.New("product price must be greater than zero")
	ErrEmptyProductName   = errors.New("product name is a required field")
	ErrInvalidProductName = fmt.Errorf("name length must be between %d and %d characters", minNameLength, maxNameLength)
)

type Product struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	Slug      slug.Slug      `json:"slug" gorm:"unique"`
	PublicID  string         `json:"public_id" gorm:"unique"`
	Name      string         `json:"name"`
	Price     float64        `json:"price"`
	Quantity  int32          `json:"quantity"`
	Enabled   bool           `json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func New(name string, price float64, quantity int32) (*Product, error) {
	if len(name) == 0 {
		return nil, ErrEmptyProductName
	}

	s, err := slug.New(name)
	if err != nil {
		return nil, err
	}

	product := Product{
		Name:     name,
		Price:    price,
		Quantity: quantity,
		PublicID: uuid.NewString(),
		Slug:     s.GetSlug(),
		Enabled:  true,
	}

	if err = product.validate(); err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) validate() error {
	invalidProdName := len(p.Name) < minNameLength || len(p.Name) >= maxNameLength

	if invalidProdName {
		return ErrInvalidProductName
	}

	if p.Quantity < 1 {
		return ErrInvalidQuantity
	}

	if p.Price < 1 {
		return ErrInvalidPrice
	}

	return nil
}
