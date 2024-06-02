package password_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/domain/value-object/password"
	"github.com/stretchr/testify/assert"
)

func TestPasswordValueObject(t *testing.T) {
	t.Run("Should create a password instance successfully", func(t *testing.T) {
		_, err := password.New("@MyPassword123")

		assert.Nil(t, err)
	})

	t.Run("Should not create a password without special character", func(t *testing.T) {
		_, err := password.New("MyPassword123")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "password must contain at least one special character")
	})

	t.Run("Should not create a password without uppercase letter", func(t *testing.T) {
		_, err := password.New("@mymassword123")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "password must contain at least one uppercase letter")
	})

	t.Run("Should not create a password without lower letter", func(t *testing.T) {
		_, err := password.New("@MYPASSWORD123")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "password must contain at least one lowercase letter")
	})

	t.Run("Should not create a password without a digit", func(t *testing.T) {
		_, err := password.New("MyPassword")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "password must contain at least one digit")
	})

	t.Run("Should not create a password without minimum character", func(t *testing.T) {
		_, err := password.New("@My1")

		assert.NotNil(t, err)
		assert.EqualError(t, err, password.ErrPasswordTooShort.Error())
	})

	t.Run("Should not create a password if you exceed the max limit", func(t *testing.T) {
		longPassword := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
		_, err := password.New(longPassword)

		assert.NotNil(t, err)
		assert.EqualError(t, err, password.ErrPasswordTooLong.Error())
	})
}
