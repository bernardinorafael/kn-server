package product

import (
	"errors"
	"fmt"
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/slug"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	minNameLength = 3
	maxNameLength = 120
)

var (
	ErrInvalidQuantity           = errors.New("product quantity cannot be zero")
	ErrManipulateDisabledProduct = errors.New("cannot manipulate a disabled product")
	ErrInvalidPrice              = errors.New("product price must be greater than zero")
	ErrEmptyProductName          = errors.New("product name is a required field")
	ErrInvalidProductName        = fmt.Errorf("name length must be between %d and %d characters", minNameLength, maxNameLength)
)

type Product struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	Slug      slug.Slug      `json:"slug" gorm:"unique"`
	PublicID  string         `json:"public_id" gorm:"unique"`
	Name      string         `json:"name"`
	Image     string         `json:"image"`
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
	if len(p.Name) < minNameLength || len(p.Name) >= maxNameLength {
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

func (p *Product) ChangePrice(price float64) error {
	if price < 1 {
		return ErrInvalidPrice
	}

	if !p.Enabled {
		return ErrManipulateDisabledProduct
	}

	p.Price = price
	return nil
}

func (p *Product) IncreaseQuantity(quantity int32) error {
	if quantity < 1 {
		return ErrInvalidQuantity
	}
	if !p.Enabled {
		return ErrManipulateDisabledProduct
	}

	p.Quantity += quantity
	return nil
}

func (p *Product) ChangeName(name string) error {
	if !p.Enabled {
		return ErrManipulateDisabledProduct
	}

	if len(name) == 0 {
		return ErrEmptyProductName
	}

	if len(p.Name) < minNameLength || len(p.Name) >= maxNameLength {
		return ErrInvalidProductName
	}

	s, err := slug.New(name)
	if err != nil {
		return err
	}

	p.Name = name
	p.Slug = s.GetSlug()

	return nil
}

func (p *Product) SetImageURL(imageURL string) {
	p.Image = imageURL
}

func (p *Product) GetStatus() bool {
	return p.Enabled
}

func (p *Product) Disable() {
	p.Enabled = false
}

func (p *Product) Enable() {
	p.Enabled = true
}
