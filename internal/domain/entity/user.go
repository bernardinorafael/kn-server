package entity

import (
	"errors"
	"regexp"

	"github.com/bernardinorafael/kn-server/helper/crypto"
	"github.com/bernardinorafael/kn-server/internal/domain/value-object/email"

	"gorm.io/gorm"
)

var (
	ErrInvalidNameLength = errors.New("name must be at least 3 characters long")
	ErrInvalidFullName   = errors.New("invalid name, must contain name and full name")
	ErrShortPassword     = errors.New("password must contain at least 6 char")
)

type User struct {
	gorm.Model

	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func NewUser(name, mail, password string) (*User, error) {
	if len(name) < 3 {
		return nil, ErrInvalidNameLength
	}

	ok, _ := regexp.MatchString("^[A-Za-zÀ-ÿ]+(?:\\s[A-Za-zÀ-ÿ]+)+$", name)
	if !ok {
		return nil, ErrInvalidFullName
	}

	address, err := email.New(mail)
	if err != nil {
		return nil, err
	}

	err = address.Validate()
	if err != nil {
		return nil, err
	}

	// TODO: consider changing password rules to improve security
	if len(password) < 6 {
		return nil, ErrShortPassword
	}

	encrypted, err := crypto.Make(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:     name,
		Email:    address.GetString(),
		Password: encrypted,
	}, nil
}
