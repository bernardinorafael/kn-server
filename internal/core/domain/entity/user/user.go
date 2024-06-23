package user

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/cpf"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/password"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	minNameLength = 3
)

var (
	ErrInvalidNameLength = fmt.Errorf("name must be at least %d characters long", minNameLength)
	ErrInvalidFullName   = errors.New("incorrect name, must contain valid first and last name")
	ErrEmptyUserName     = errors.New("name is a required field")
)

type User struct {
	ID        int               `json:"id" gorm:"primaryKey"`
	Name      string            `json:"name"`
	Email     email.Email       `json:"email" gorm:"unique"`
	PublicID  string            `json:"public_id" gorm:"unique"`
	Document  cpf.CPF           `json:"document" gorm:"unique"`
	Enabled   bool              `json:"enabled"`
	Password  password.Password `json:"password"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `json:"deleted_at"`
}

func New(userName, userEmail, userPassword, userDoc string) (*User, error) {
	address, err := email.New(userEmail)
	if err != nil {
		return nil, err
	}

	p, err := password.New(userPassword)
	if err != nil {
		return nil, err
	}

	hashed, err := p.ToEncrypted()
	if err != nil {
		return nil, err
	}

	document, err := cpf.New(userDoc)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     userName,
		PublicID: uuid.NewString(),
		Email:    address.ToEmail(),
		Document: document.ToCPF(),
		Password: hashed,
		Enabled:  false,
	}

	if err = user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) validate() error {
	if len(u.Name) == 0 {
		return ErrEmptyUserName
	}

	if len(u.Name) < minNameLength {
		return ErrInvalidNameLength
	}

	u.Name = strings.TrimSpace(u.Name)

	fullNamePattern := "^[A-Za-zÀ-ÿ]+(?:\\s[A-Za-zÀ-ÿ]+)+$"
	matched, _ := regexp.MatchString(fullNamePattern, u.Name)
	if !matched {
		return ErrInvalidFullName
	}

	return nil
}
