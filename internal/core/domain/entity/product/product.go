package product

import (
	"errors"
	"fmt"
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/slug"
	"github.com/google/uuid"
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
	PublicID  string    `json:"public_id"`
	Slug      slug.Slug `json:"slug"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

func New(name string, price float64, quantity int) (*Product, error) {
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
	if len(p.Name) > maxNameLength {
		return ErrInvalidProductName
	}

	if len(p.Name) < minNameLength {
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

func (p *Product) IncreaseQuantity(quantity int) error {
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

func (p *Product) SetImage(url string) {
	p.Image = url
}

func (p *Product) GetStatus() bool {
	return p.Enabled
}

func (p *Product) ChangeStatus(status bool) {
	p.Enabled = status
}
