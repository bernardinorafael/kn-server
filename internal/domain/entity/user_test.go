package entity_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/domain/entity"

	"github.com/stretchr/testify/assert"
)

func TestUser_New(t *testing.T) {
	name := "john doe"
	email := "john_doe@gmail.com"
	password := "@MyPassword123"

	t.Run("Should create an entity", func(t *testing.T) {
		_, err := entity.NewUser(name, email, password)
		assert.Nil(t, err)
	})

	t.Run("Should entity name have at least 3 char", func(t *testing.T) {
		_, err := entity.NewUser("jo", email, password)
		assert.EqualError(t, err, "name must be at least 3 characters long")
	})

	t.Run("Should entity have full name", func(t *testing.T) {
		_, err := entity.NewUser("john", email, password)
		assert.EqualError(t, err, "invalid name, must contain name and full name")

		_, err = entity.NewUser("john ", email, password)
		assert.EqualError(t, err, "invalid name, must contain name and full name")

		_, err = entity.NewUser("john doe ", email, password)
		assert.EqualError(t, err, "invalid name, must contain name and full name")
	})
}
