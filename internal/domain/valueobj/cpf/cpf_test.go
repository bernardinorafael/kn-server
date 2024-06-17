package cpf_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/cpf"
	"github.com/stretchr/testify/assert"
)

func TestCPFValueObject_New(t *testing.T) {
	t.Run("Should validate CPF digit", func(t *testing.T) {
		valid := "42008790002"
		document, err := cpf.New(valid)

		assert.Nil(t, err)
		assert.Equal(t, document.ToCPF(), cpf.CPF(valid))
	})

	t.Run("Should throw error if digit validator is wrong", func(t *testing.T) {
		invalid := "12028730002"
		_, err := cpf.New(invalid)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid cpf")
	})

	t.Run("Should throw error if document do not have exact 11 digits", func(t *testing.T) {
		invalid := "1202873"
		_, err := cpf.New(invalid)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "cpf must have 11 characters")
	})
}
