package user_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserEntity(t *testing.T) {
	t.Run("Should create a new user instance successfully", func(t *testing.T) {
		u, err := user.New("john doe", "john_doe@email.com", "@Password123")

		assert.Nil(t, err)
		assert.Equal(t, u.Name, "john doe")
		assert.Equal(t, string(u.Email), "john_doe@email.com")
		assert.NotEqual(t, u.Password, "@Password123")

		_, err = uuid.Parse(u.PublicID)
		assert.Nil(t, err)
	})

	t.Run("Should validate user name length", func(t *testing.T) {
		_, err := user.New("jo", "john_doe@email.com", "@Password123")

		assert.NotNil(t, err)
		assert.EqualError(t, err, user.ErrInvalidNameLength.Error())
	})

	t.Run("Should user name have first and last name", func(t *testing.T) {
		_, err := user.New("joe", "john_doe@email.com", "@Password123")

		assert.NotNil(t, err)
		assert.EqualError(t, err, user.ErrInvalidFullName.Error())
	})
}