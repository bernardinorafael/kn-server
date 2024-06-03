package entity

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/slug"
	"gorm.io/gorm"
)

var (
	ErrShortProductName = errors.New("product name must be at least 3 characters long")
	ErrInvalidQuantity  = errors.New("product quantity cannot be less than 1")
)

type Product struct {
	gorm.Model

	Slug     slug.Slug `json:"slug" gorm:"unique"`
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
	Quantity int32     `json:"quantity"`
	Status   bool      `json:"status" gorm:"default:true"`
}

func NewProduct(name string, price float64, quantity int32) (*Product, error) {
	if len(name) <= 3 {
		return nil, ErrShortProductName
	}

	if quantity < 1 {
		return nil, ErrInvalidQuantity
	}

	s, err := slug.New(name)
	if err != nil {
		return nil, err
	}

	return &Product{
		Name:     name,
		Price:    price,
		Quantity: quantity,
		Slug:     s.GetSlug(),
	}, nil
}
