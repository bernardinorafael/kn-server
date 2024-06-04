package product_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"
	"github.com/stretchr/testify/assert"
)

func TestProductEntity(t *testing.T) {
	t.Run("Should create a new product instance", func(t *testing.T) {
		p, err := product.New("my product name", 300.1, 100)

		assert.Nil(t, err)
		assert.Equal(t, string(p.Slug), "my-product-name")
	})

	t.Run("Should throw an error if name is less than 3 characters", func(t *testing.T) {
		_, err := product.New("pr", 300.1, 100)

		assert.NotNil(t, err)
		assert.EqualError(t, err, product.ErrInvalidProductName.Error())
	})

	t.Run("Should throw an error if name is greater than 120 characters", func(t *testing.T) {
		_, err := product.New(
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque cursus at sapien id pretium. Mauris convallis, urna eget.",
			300.1,
			100,
		)

		assert.NotNil(t, err)
		assert.EqualError(t, err, product.ErrInvalidProductName.Error())
	})

	t.Run("Should product price be a valid integer", func(t *testing.T) {
		_, err := product.New("my product name", 0, 100)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "product price must be greater than zero")
	})

	t.Run("Should product quantity be a valid integer", func(t *testing.T) {
		_, err := product.New("my product name", 300.1, 0)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "product quantity cannot be zero")
	})
}
