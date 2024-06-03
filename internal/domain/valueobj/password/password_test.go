package password_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/password"
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

	t.Run("Should password be encrypted correctly", func(t *testing.T) {
		pass, _ := password.New("@MyPassword123")
		encrypted, err := pass.ToEncrypted()

		assert.Nil(t, err)
		assert.NotEqual(t, encrypted, "@MyPassword123")
	})

	t.Run("Should compare passwords correctly", func(t *testing.T) {
		encryptedPassword := "$2a$10$mTrfubJ5HjIci00eP/fCzuZOe/2YOYP9PGWLOh6y/E/YCtvRlOylO"

		pass, err := password.New("@MyPassword123")
		matched := pass.Compare(password.Password(encryptedPassword))

		assert.Nil(t, err)
		assert.Nil(t, matched)
	})

	t.Run("Should throw an error if compare fails", func(t *testing.T) {
		wrongEncrypted := "$2a$70$mTrfubJ5HjIci00eP/fCzuZOe/21OYP9PGWLOh6y/E/YCtvR3lOylO"

		pass, err := password.New("@MyPassword123")
		matched := pass.Compare(password.Password(wrongEncrypted))

		assert.Nil(t, err)
		assert.NotNil(t, matched)
		assert.EqualError(t, matched, "provided password does not match")
	})
}
