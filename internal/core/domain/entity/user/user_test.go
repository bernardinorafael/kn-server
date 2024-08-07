package user_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/password"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserEntity_New(t *testing.T) {
	t.Run("Should throw an error if ID is not a valid uuid", func(t *testing.T) {
		u, err := user.New(user.Params{
			PublicID: "invalid-id",
			Name:     "john doe",
			Email:    "john_doe@email.com",
			Password: "@Password123",
			Phone:    "11978761232",
			TeamID:   nil,
		})

		assert.EqualError(t, err, "invalid id, must be a valid uuid format")
		assert.Nil(t, u)
	})

	t.Run("Should create a new user instance successfully", func(t *testing.T) {
		u, err := user.New(user.Params{
			PublicID: uuid.NewString(),
			Name:     "john doe",
			Email:    "john_doe@email.com",
			Password: "@Password123",
			Phone:    "11978761232",
			TeamID:   nil,
		})
		assert.Nil(t, err)

		err = u.EncryptPassword()
		assert.Nil(t, err)

		assert.Equal(t, u.Name(), "john doe")
		assert.Equal(t, u.Email(), email.Email("john_doe@email.com"))
		assert.NotEqual(t, u.Password(), password.Password("@Password123"))
	})

	t.Run("Should validate user name length", func(t *testing.T) {
		_, err := user.New(user.Params{
			PublicID: uuid.NewString(),
			Name:     "jo",
			Email:    "john_doe@email.com",
			Password: "@Password123",
			Phone:    "11978761232",
			TeamID:   nil,
		})

		assert.NotNil(t, err)
		assert.EqualError(t, err, user.ErrInvalidNameLength.Error())
	})

	t.Run("Should user name have first and last name", func(t *testing.T) {
		_, err := user.New(user.Params{
			PublicID: uuid.NewString(),
			Name:     "joe",
			Email:    "john_doe@email.com",
			Password: "@Password123",
			Phone:    "11978761232",
			TeamID:   nil,
		})

		assert.NotNil(t, err)
		assert.EqualError(t, err, user.ErrInvalidFullName.Error())
	})
}

func TestUser_ChangeStatus(t *testing.T) {
	t.Run("Should only change a status by a valid one", func(t *testing.T) {
		u, _ := user.New(user.Params{
			PublicID: uuid.NewString(),
			Name:     "joe doe",
			Email:    "john_doe@email.com",
			Password: "@Password123",
			Phone:    "11978761232",
			TeamID:   nil,
		})

		err := u.ChangeStatus(user.StatusEnabled)
		assert.Nil(t, err)
	})

	t.Run("Should not change status by a invalid one", func(t *testing.T) {
		u, _ := user.New(user.Params{
			PublicID: uuid.NewString(),
			Name:     "joe doe",
			Email:    "john_doe@email.com",
			Password: "@Password123",
			Phone:    "11978761232",
			TeamID:   nil,
		})

		err := u.ChangeStatus(user.Status(99))
		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid user status")
	})

	t.Run("Should return status string from user", func(t *testing.T) {
		u, _ := user.New(user.Params{
			PublicID: uuid.NewString(),
			Name:     "joe doe",
			Email:    "john_doe@email.com",
			Password: "@Password123",
			Phone:    "11978761232",
			TeamID:   nil,
		})

		assert.Equal(t, u.StatusString(), "pending")
	})
}
