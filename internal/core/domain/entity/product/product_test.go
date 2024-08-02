package product_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductEntity_New(t *testing.T) {
	t.Run("Should create a new product instance", func(t *testing.T) {
		p, err := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300.1,
			Quantity: 100,
			Enabled:  false,
		})
		assert.Nil(t, err)
		assert.Equal(t, string(p.Slug()), "my-product-name")
	})

	t.Run("Should throw an error if name is less than 3 characters", func(t *testing.T) {
		_, err := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "pr",
			Image:    "image-url-test",
			Price:    300.1,
			Quantity: 100,
			Enabled:  false,
		})
		assert.NotNil(t, err)
		assert.EqualError(t, err, product.ErrInvalidProductName.Error())
	})

	t.Run("Should throw an error if name is greater than 120 characters", func(t *testing.T) {
		_, err := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque cursus at sapien id pretium. Mauris convallis, urna eget.",
			Image:    "image-url-test",
			Price:    300.1,
			Quantity: 100,
			Enabled:  false,
		})
		assert.NotNil(t, err)
		assert.EqualError(t, err, product.ErrInvalidProductName.Error())
	})

	t.Run("Should product price be a valid integer", func(t *testing.T) {
		_, err := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    0,
			Quantity: 10,
			Enabled:  false,
		})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "product price must be greater than zero")
	})
}

func TestProductEntity_ChangeStatus(t *testing.T) {
	t.Run("Should disable a product", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300.1,
			Quantity: 100,
			Enabled:  true,
		})

		p.ChangeStatus(false)
		assert.False(t, p.Enabled())
	})

	t.Run("Should enable a product", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300.1,
			Quantity: 100,
			Enabled:  false,
		})

		p.ChangeStatus(true)
		assert.True(t, p.Enabled())
	})
}

func TestProductEntity_IncreasePrice(t *testing.T) {
	t.Run("Should be possible to increase product price", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300,
			Quantity: 100,
			Enabled:  false,
		})

		err := p.ChangePrice(100)
		assert.Nil(t, err)
		assert.Equal(t, p.Price(), float64(100))
	})

	t.Run("Should not be able to increase price if the inc number is lesser than zero", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300,
			Quantity: 100,
			Enabled:  false,
		})

		err := p.ChangePrice(0)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "product price must be greater than zero")
	})

	t.Run("Should not be possible to increase price if product is disabled", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300,
			Quantity: 100,
			Enabled:  false,
		})
		p.ChangeStatus(false)

		err := p.ChangePrice(10)
		assert.NotNil(t, err)

		assert.EqualError(t, err, "cannot manipulate a disabled product")
	})
}

func TestProductEntity_IncQuantity(t *testing.T) {
	t.Run("Should increment product quantity", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300,
			Quantity: 10,
			Enabled:  false,
		})

		err := p.IncreaseQuantity(10)

		assert.Nil(t, err)
		assert.Equal(t, p.Quantity(), int32(20))
	})

	t.Run("Should not be able to inc a product quantity with zero value", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300,
			Quantity: 10,
			Enabled:  false,
		})

		err := p.IncreaseQuantity(0)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "product quantity cannot be zero")
	})

	t.Run("Should not be possible to increase quantity if product is disabled", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300,
			Quantity: 10,
			Enabled:  false,
		})
		p.ChangeStatus(false)

		err := p.IncreaseQuantity(10)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "cannot manipulate a disabled product")
	})
}

func TestProductEntity_ChangeName(t *testing.T) {
	t.Run("Should change product name", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300,
			Quantity: 10,
			Enabled:  false,
		})

		err := p.ChangeName("other product name")
		assert.Nil(t, err)
		assert.Equal(t, p.Name(), "other product name")
		assert.Equal(t, string(p.Slug()), "other-product-name")
	})

	t.Run("Should not possible to change name of a disabled product", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300,
			Quantity: 10,
			Enabled:  true,
		})
		p.ChangeStatus(false)

		err := p.ChangeName("other product name")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "cannot manipulate a disabled product")
		assert.Equal(t, p.Name(), "my product name")
	})

	t.Run("Should not possible to change the name is name attribute is empty", func(t *testing.T) {
		p, _ := product.New(product.Params{
			PublicID: uuid.NewString(),
			Name:     "my product name",
			Image:    "image-url-test",
			Price:    300,
			Quantity: 10,
			Enabled:  true,
		})

		err := p.ChangeName("")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "product name is a required field")
		assert.Equal(t, p.Name(), "my product name")
	})
}
