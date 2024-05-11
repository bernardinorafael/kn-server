package entity_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/domain/entity"

	"github.com/stretchr/testify/assert"
)

func TestUser_New(t *testing.T) {
	name := "john doe"
	email := "john_doe@gmail.com"
	password := "abcd1234"

	t.Run("should create an entity", func(t *testing.T) {
		u, err := entity.NewUser(name, email, password)

		assert.Nil(t, err)
		assert.Equal(t, name, u.Name)
		assert.Equal(t, email, u.Email)
		assert.NotEqual(t, password, u.Password)
	})

	t.Run("should not create an entity if email is invalid", func(t *testing.T) {
		_, err := entity.NewUser(name, "john_doe@nothing", password)
		assert.EqualError(t, err, "invalid email address format")

		_, err = entity.NewUser(name, "john", password)
		assert.EqualError(t, err, "invalid email address format")

		_, err = entity.NewUser(name, "john@gmail", password)
		assert.EqualError(t, err, "invalid email address format")

		_, err = entity.NewUser(name, "john@gmail.", password)
		assert.EqualError(t, err, "invalid email address format")

		_, err = entity.NewUser(name, "", password)
		assert.EqualError(t, err, "invalid email address format")
	})

	t.Run("should entity name have at least 3 char", func(t *testing.T) {
		_, err := entity.NewUser("jo", email, password)
		assert.EqualError(t, err, "name must be at least 3 characters long")
	})

	t.Run("should entity have full name", func(t *testing.T) {
		_, err := entity.NewUser("john", email, password)
		assert.EqualError(t, err, "invalid name, must contain name and full name")

		_, err = entity.NewUser("john ", email, password)
		assert.EqualError(t, err, "invalid name, must contain name and full name")

		_, err = entity.NewUser("john doe ", email, password)
		assert.EqualError(t, err, "invalid name, must contain name and full name")
	})
}
