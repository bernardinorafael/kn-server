package entity

import (
	"errors"
	"regexp"

	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/password"

	"gorm.io/gorm"
)

var (
	ErrInvalidNameLength = errors.New("name must be at least 3 characters long")
	ErrInvalidFullName   = errors.New("invalid name, must contain name and full name")
	ErrShortPassword     = errors.New("password must contain at least 6 char")
)

type User struct {
	gorm.Model

	Name     string                     `json:"name"`
	Email    email.Email                `json:"email" gorm:"unique"`
	Password password.EncryptedPassword `json:"password"`
}

func NewUser(name, e, p string) (*User, error) {
	if len(name) < 3 {
		return nil, ErrInvalidNameLength
	}

	ok, _ := regexp.MatchString("^[A-Za-zÀ-ÿ]+(?:\\s[A-Za-zÀ-ÿ]+)+$", name)
	if !ok {
		return nil, ErrInvalidFullName
	}

	address, err := email.New(e)
	if err != nil {
		return nil, err
	}

	pass, err := password.New(p)
	if err != nil {
		return nil, err
	}

	encrypted, err := pass.ToEncrypted()
	if err != nil {
		return nil, err
	}

	return &User{
		Name:     name,
		Email:    address.ToEmail(),
		Password: encrypted,
	}, nil
}
