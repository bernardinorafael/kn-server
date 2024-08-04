package user

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

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
	ErrInvalidUUIDFormat = errors.New("invalid id, must be a valid uuid format")
)

// Params contains the parameters required to create a new User entity
type Params struct {
	PublicID string
	Name     string
	Email    string
	Password string
	Phone    string
	TeamID   *string
}

type User struct {
	publicID  string
	name      string
	email     email.Email
	phone     phone.Phone
	enabled   bool
	teamID    *string
	password  password.Password
	createdAt time.Time
}

func New(u Params) (*User, error) {
	address, err := email.New(u.Email)
	if err != nil {
		return nil, err
	}

	ph, err := phone.New(u.Phone)
	if err != nil {
		return nil, err
	}

	user := User{
		publicID:  u.PublicID,
		name:      u.Name,
		email:     address.ToEmail(),
		phone:     ph.Phone(),
		password:  password.Password(u.Password),
		teamID:    u.TeamID,
		enabled:   false,
		createdAt: time.Now(),
	}

	if err = user.validate(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) validate() error {
	if u.name == "" {
		return ErrEmptyUserName
	}

	if len(u.name) < minNameLength {
		return ErrInvalidNameLength
	}

	u.name = strings.TrimSpace(u.name)

	fullNamePattern := "^[A-Za-zÀ-ÿ]+(?:\\s[A-Za-zÀ-ÿ]+)+$"
	matched, _ := regexp.MatchString(fullNamePattern, u.name)
	if !matched {
		return ErrInvalidFullName
	}

	_, err := uuid.Parse(u.publicID)
	if err != nil {
		return ErrInvalidUUIDFormat
	}

	return nil
}

func (u *User) EncryptPassword() error {
	p, err := password.New(string(u.password))
	if err != nil {
		return err
	}

	hashed, err := p.ToEncrypted()
	if err != nil {
		return err
	}

	u.password = hashed
	return nil
}

func (u *User) ChangeName(newName string) error {
	u.name = newName
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

	u.phone = p.Phone()
	return nil
}

func (u *User) ChangeEmail(newEmail string) error {
	address, err := email.New(newEmail)
	if err != nil {
		return err
	}

	u.email = address.ToEmail()
	return nil
}

func (u *User) ChangeStatus(status bool) {
	u.enabled = status
}

func (u *User) PublicID() string            { return u.publicID }
func (u *User) Name() string                { return u.name }
func (u *User) Email() email.Email          { return u.email }
func (u *User) Phone() phone.Phone          { return u.phone }
func (u *User) Enabled() bool               { return u.enabled }
func (u *User) TeamID() *string             { return u.teamID }
func (u *User) Password() password.Password { return u.password }
func (u *User) CreatedAt() time.Time        { return u.createdAt }
