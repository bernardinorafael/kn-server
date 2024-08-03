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

// Params contains the parameters required to create a new Product entity
type Params struct {
	PublicID string
	Name     string
	Image    string
	Price    float64
	Quantity int
	Enabled  bool
}

type Product struct {
	publicID  string
	slug      slug.Slug
	name      string
	image     string
	price     float64
	quantity  int
	enabled   bool
	createdAt time.Time
}

func New(p Params) (*Product, error) {
	s, err := slug.New(p.Name)
	if err != nil {
		return nil, err
	}

	product := Product{
		publicID:  uuid.NewString(),
		slug:      s.GetSlug(),
		name:      p.Name,
		image:     p.Image,
		price:     p.Price,
		quantity:  p.Quantity,
		enabled:   p.Enabled,
		createdAt: time.Now(),
	}

	if err = product.validate(); err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *Product) validate() error {
	if p.Name() == "" {
		return ErrEmptyProductName
	}

	if len(p.Name()) > maxNameLength {
		return ErrInvalidProductName
	}

	if len(p.Name()) < minNameLength {
		return ErrInvalidProductName
	}

	if p.Quantity() < 1 {
		return ErrInvalidQuantity
	}

	if p.Price() < 1 {
		return ErrInvalidPrice
	}
	return nil
}

func (p *Product) ChangePrice(price float64) error {
	if price < 1 {
		return ErrInvalidPrice
	}

	if !p.Enabled() {
		return ErrManipulateDisabledProduct
	}

	p.price = price
	return nil
}

func (p *Product) IncreaseQuantity(quantity int) error {
	if quantity < 1 {
		return ErrInvalidQuantity
	}
	if !p.Enabled() {
		return ErrManipulateDisabledProduct
	}

	p.quantity += quantity
	return nil
}

func (p *Product) ChangeName(name string) error {
	if !p.Enabled() {
		return ErrManipulateDisabledProduct
	}

	if len(name) == 0 {
		return ErrEmptyProductName
	}

	if len(p.Name()) < minNameLength || len(p.Name()) >= maxNameLength {
		return ErrInvalidProductName
	}

	s, err := slug.New(name)
	if err != nil {
		return err
	}

	p.name = name
	p.slug = s.GetSlug()

	return nil
}

func (p *Product) ChangeStatus(status bool) {
	p.enabled = status
}

func (p *Product) PublicID() string     { return p.publicID }
func (p *Product) Slug() slug.Slug      { return p.slug }
func (p *Product) Name() string         { return p.name }
func (p *Product) Image() string        { return p.image }
func (p *Product) Price() float64       { return p.price }
func (p *Product) Quantity() int        { return p.quantity }
func (p *Product) Enabled() bool        { return p.enabled }
func (p *Product) CreatedAt() time.Time { return p.createdAt }
