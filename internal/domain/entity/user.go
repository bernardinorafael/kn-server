package entity

import (
	"errors"
	"regexp"

	"github.com/bernardinorafael/kn-server/helper/crypto"
	"gorm.io/gorm"
)

var (
	ErrInvalidNameLength   = errors.New("name must be at least 3 characters long")
	ErrInvalidFullName     = errors.New("invalid name, must contain name and full name")
	ErrInvalidEmailAddress = errors.New("invalid email address format")
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(name, email, password string) (*User, error) {
	if len(name) <= 3 {
		return nil, ErrInvalidNameLength
	}

	ok, _ := regexp.MatchString("^[A-Za-zÀ-ÿ]+(?:\\s[A-Za-zÀ-ÿ]+)+$", name)
	if !ok {
		return nil, ErrInvalidFullName
	}

	ok, _ = regexp.MatchString("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", email)
	if !ok {
		return nil, ErrInvalidEmailAddress
	}

	encrypted, err := crypto.Make(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     name,
		Email:    email,
		Password: encrypted,
	}

	return user, nil
}
