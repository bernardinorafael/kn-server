package password

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLength = 6
	maxPasswordLength = 64
)

var (
	ErrPasswordTooShort       = fmt.Errorf("password must be at least %d characters long", minPasswordLength)
	ErrPasswordTooLong        = fmt.Errorf("password must not exceed %d characters", maxPasswordLength)
	ErrMissingSpecialChar     = errors.New("password must contain at least one special character")
	ErrMissingUppercaseLetter = errors.New("password must contain at least one uppercase letter")
	ErrMissingLowercaseLetter = errors.New("password must contain at least one lowercase letter")
	ErrPasswordDoesNotMatch   = errors.New("provided password does not match")
	ErrMissingDigit           = errors.New("password must contain at least one digit")
)

type Password string

type password struct {
	password string
}

func New(rawPassword string) (*password, error) {
	p := password{password: rawPassword}

	err := p.validate()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *password) validate() error {
	if len(p.password) == 0 {
		return errors.New("password is a required field")
	}

	if len(p.password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	if len(p.password) > maxPasswordLength {
		return ErrPasswordTooLong
	}

	matched, _ := regexp.MatchString(`[0-9]`, p.password)
	if !matched {
		return ErrMissingDigit
	}

	matched, _ = regexp.MatchString(`[a-z]`, p.password)
	if !matched {
		return ErrMissingLowercaseLetter
	}

	matched, _ = regexp.MatchString(`[A-Z]`, p.password)
	if !matched {
		return ErrMissingUppercaseLetter
	}

	specialCharPattern := `[!@#~$%^&*()_+|{}\[\]:;"'<>,.?/\-]`
	matched, _ = regexp.MatchString(specialCharPattern, p.password)
	if !matched {
		return ErrMissingSpecialChar
	}

	return nil
}

func (p *password) Encrypt() (Password, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(p.password), 10)
	if err != nil {
		return "", err
	}
	return Password(encrypted), nil
}

func (p *password) Compare(hashed Password, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return ErrPasswordDoesNotMatch
	}
	return nil
}
