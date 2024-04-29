package entity

import (
	"errors"
	"regexp"
	"time"

	"github.com/bernardinorafael/kn-server/helper/crypto"
	"github.com/google/uuid"
)

var (
	ErrInvalidNameLength   = errors.New("name must be at least 3 characters long")
	ErrInvalidFullName     = errors.New("invalid name, must contain name and full name")
	ErrInvalidEmailAddress = errors.New("invalid email address format")
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func New(name, email, password string) (*User, error) {
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
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  encrypted,
		CreatedAt: time.Now(),
	}

	return user, nil
}
