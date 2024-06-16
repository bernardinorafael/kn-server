package user

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/password"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	minNameLength = 3
)

var (
	ErrInvalidNameLength = fmt.Errorf("name must be at least %d characters long", minNameLength)
	ErrInvalidFullName   = errors.New("name must contain both first name and last name")
	ErrShortPassword     = errors.New("password must be at least 6 characters long")
)

type User struct {
	ID        int               `json:"id" gorm:"primaryKey"`
	Name      string            `json:"name"`
	Email     email.Email       `json:"email" gorm:"unique"`
	PublicID  string            `json:"public_id" gorm:"unique"`
	Enabled   bool              `json:"enabled"`
	Password  password.Password `json:"password"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `json:"deleted_at"`
}

func New(userName, userEmail, userPass string) (*User, error) {
	address, err := email.New(userEmail)
	if err != nil {
		return nil, err
	}

	passw, err := password.New(userPass)
	if err != nil {
		return nil, err
	}

	encrypted, err := passw.ToEncrypted()
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     userName,
		PublicID: uuid.NewString(),
		Email:    address.ToEmail(),
		Password: encrypted,
		Enabled:  false,
	}

	if err = user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) validate() error {
	if len(u.Name) < minNameLength {
		return ErrInvalidNameLength
	}

	fullNamePattern := "^[A-Za-zÀ-ÿ]+(?:\\s[A-Za-zÀ-ÿ]+)+$"
	matched, _ := regexp.MatchString(fullNamePattern, u.Name)
	if !matched {
		return ErrInvalidFullName
	}

	return nil
}
