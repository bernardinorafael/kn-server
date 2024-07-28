package phone_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/phone"
	"github.com/stretchr/testify/assert"
)

func TestPhone_New(t *testing.T) {
	t.Run("Should validate a phone", func(t *testing.T) {
		p, err := phone.New("11988091232")

		assert.Nil(t, err)
		assert.Equal(t, p.ToPhone(), phone.Phone("11988091232"))
	})

	t.Run("Should return complete string phone", func(t *testing.T) {
		p, err := phone.New("11988091232")

		phoneStr := p.ToPhone()

		assert.Nil(t, err)
		assert.Equal(t, phoneStr, phone.Phone("11988091232"))
	})

	t.Run("Should phone contains actually 11 length size", func(t *testing.T) {
		_, err := phone.New("119889122")

		assert.NotNil(t, err)
		assert.EqualError(t, err, err.Error())
	})

	t.Run("Should throw an error if phone had an invalid area code", func(t *testing.T) {
		_, err := phone.New("09988566239")

		assert.NotNil(t, err)
		assert.EqualError(t, err, err.Error())
	})

	t.Run("Should throw error if local part does not init with 9", func(t *testing.T) {
		_, err := phone.New("88388063331")

		assert.NotNil(t, err)
		assert.EqualError(t, err, err.Error())
	})
}
