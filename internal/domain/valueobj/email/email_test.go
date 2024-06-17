package email_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/email"
	"github.com/stretchr/testify/assert"
)

func TestEmailValueObject(t *testing.T) {
	t.Run("Should create the Address format", func(t *testing.T) {
		address, err := email.New("john_doe@email.com")

		assert.Nil(t, err)
		assert.Equal(t, email.Email("john_doe"), address.GetLocalPart())
		assert.Equal(t, email.Email("email.com"), address.GetDomainPart())
	})

	t.Run("Should not create an email with invalid characters", func(t *testing.T) {
		_, err := email.New("**!) !&#%^&*()+=[]{}|;:'\",<>?/\\~`@email.com")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid email address")
	})

	t.Run("Should local part have at least 3 characters", func(t *testing.T) {
		_, err := email.New("jo@email.com")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "local part of the email must have at least 3 characters")
	})

	t.Run("Should not local part have two period in sequence", func(t *testing.T) {
		_, err := email.New("john..doe@email.com")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid email address")
	})

	t.Run("Should domain part have at least 3 characters", func(t *testing.T) {
		_, err := email.New("john_doe@em.com")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "domain part email must have at least 3 characters")
	})

	t.Run("Should method GetString return the email correctly", func(t *testing.T) {
		address, err := email.New("john_doe@email.com")

		email := address.ToEmail()

		assert.Nil(t, err)
		assert.Contains(t, string(email), "@")
	})

	t.Run("Should method GetDomainPart return the email correctly", func(t *testing.T) {
		address, err := email.New("john_doe@email.com")

		domainPart := address.GetDomainPart()

		assert.Nil(t, err)
		assert.NotContains(t, domainPart, "@")
		assert.Contains(t, domainPart, ".")
		assert.Equal(t, domainPart, email.Email("email.com"))
	})

	t.Run("Should method GetLocalPart return the email correctly", func(t *testing.T) {
		address, err := email.New("john_doe@email.com")

		localPart := address.GetLocalPart()

		assert.Nil(t, err)
		assert.NotContains(t, string(localPart), "@")
		assert.Equal(t, localPart, email.Email("john_doe"))
	})

	t.Run("Should DomainPart ends with . plus domain", func(t *testing.T) {
		_, err := email.New("john_doe@email")

		assert.NotNil(t, err)
	})
}
