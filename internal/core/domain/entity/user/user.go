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
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/phone"
	"github.com/google/uuid"
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
	PublicID  string            `json:"public_id"`
	Name      string            `json:"name"`
	Email     email.Email       `json:"email"`
	Document  cpf.CPF           `json:"document"`
	Phone     phone.Phone       `json:"phone"`
	Enabled   bool              `json:"enabled"`
	Password  password.Password `json:"password"`
	CreatedAt time.Time         `json:"created_at"`
}

func New(newName, newEmail, newPass, newDoc, newPhone string) (*User, error) {
	address, err := email.New(newEmail)
	if err != nil {
		return nil, err
	}

	ph, err := phone.New(newPhone)
	if err != nil {
		return nil, err
	}

	p, err := password.New(newPass)
	if err != nil {
		return nil, err
	}

	hashed, err := p.ToEncrypted()
	if err != nil {
		return nil, err
	}

	document, err := cpf.New(newDoc)
	if err != nil {
		return nil, err
	}

	user := User{
		Name:     newName,
		PublicID: uuid.NewString(),
		Email:    address.ToEmail(),
		Document: document.ToCPF(),
		Phone:    ph.ToPhone(),
		Password: hashed,
		Enabled:  false,
	}

	if err = user.validate(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) validate() error {
	if u.Name == "" {
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

func (u *User) ChangeName(newName string) error {
	u.Name = newName
	if err := u.validate(); err != nil {
		return err
	}
	return nil
}

func (u *User) ChangePhone(newPhone string) error {
	p, err := phone.New(newPhone)
	if err != nil {
		return err
	}

	u.Phone = p.ToPhone()
	return nil
}

func (u *User) ChangeDocument(newDocument string) error {
	doc, err := cpf.New(newDocument)
	if err != nil {
		return err
	}

	u.Document = doc.ToCPF()
	return nil
}

func (u *User) ChangeEmail(newEmail string) error {
	address, err := email.New(newEmail)
	if err != nil {
		return err
	}

	u.Email = address.ToEmail()
	return nil
}
