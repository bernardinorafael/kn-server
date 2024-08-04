package money_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/money"
	"github.com/stretchr/testify/assert"
)

func TestMoney_New(t *testing.T) {
	t.Run("Should init Money correctly", func(t *testing.T) {
		_, err := money.New(100)
		assert.Nil(t, err)
	})

	t.Run("Should receive an error if get a negative amount", func(t *testing.T) {
		_, err := money.New(0)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "monetary amounts cannot be negative")
	})

	t.Run("Should receive an error if the max value exceeded", func(t *testing.T) {
		_, err := money.New(9999999999999999)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "maximum amount value has been exceeded")
	})
}
