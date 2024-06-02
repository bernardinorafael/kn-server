package email_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/domain/value-object/email"
	"github.com/stretchr/testify/assert"
)

func TestEmailValueObject(t *testing.T) {
	t.Run("Should create the Address format", func(t *testing.T) {
		address, err := email.New("john_doe@email.com")

		assert.Nil(t, err)
		assert.Equal(t, "john_doe", address.GetLocalPart())
		assert.Equal(t, "email.com", address.GetDomainPart())
	})

	t.Run("Should not create an email with invalid characters", func(t *testing.T) {
		address, _ := email.New("**!) !&#%^&*()+=[]{}|;:'\",<>?/\\~`@email.com")
		err := address.Validate()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "Email has some invalid special characters")
	})

	t.Run("Should local part have at least 3 characters", func(t *testing.T) {
		address, _ := email.New("jo@email.com")
		err := address.Validate()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "Local part of the email must have at least 3 characters")
	})

	t.Run("Should not local part have two period in sequence", func(t *testing.T) {
		address, _ := email.New("john..doe@email.com")
		err := address.Validate()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "Invalid email address")
	})

	t.Run("Should domain part have at least 3 characters", func(t *testing.T) {
		address, _ := email.New("john_doe@em.com")
		err := address.Validate()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "Domain part email must have at least 3 characters")
	})

	t.Run("Should method GetString return the email correctly", func(t *testing.T) {
		address, _ := email.New("john_doe@email.com")

		err := address.Validate()
		emailString := address.GetString()

		assert.Nil(t, err)
		assert.Contains(t, emailString, "@")
	})

	t.Run("Should method GetDomainPart return the email correctly", func(t *testing.T) {
		address, _ := email.New("john_doe@email.com")

		err := address.Validate()
		emailString := address.GetDomainPart()

		assert.Nil(t, err)
		assert.NotContains(t, emailString, "@")
		assert.Contains(t, emailString, ".")
		assert.Equal(t, emailString, "email.com")
	})

	t.Run("Should method GetLocalPart return the email correctly", func(t *testing.T) {
		address, _ := email.New("john_doe@email.com")

		err := address.Validate()
		emailString := address.GetLocalPart()

		assert.Nil(t, err)
		assert.NotContains(t, emailString, "@")
		assert.Equal(t, emailString, "john_doe")
	})
}
