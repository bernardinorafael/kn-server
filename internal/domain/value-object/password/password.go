package password

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	minPasswordLength = 6
	maxPasswordLength = 64
)

var (
	ErrEmptyPassword          = errors.New("password cannot be empty")
	ErrPasswordTooShort       = fmt.Errorf("password must be at least %d characters long", minPasswordLength)
	ErrPasswordTooLong        = fmt.Errorf("password must not exceed %d characters", maxPasswordLength)
	ErrMissingSpecialChar     = errors.New("password must contain at least one special character")
	ErrMissingUppercaseLetter = errors.New("password must contain at least one uppercase letter")
	ErrMissingLowercaseLetter = errors.New("password must contain at least one lowercase letter")
	ErrMissingDigit           = errors.New("password must contain at least one digit")
)

type Password struct {
	password string
}

func New(rawPassword string) (*Password, error) {
	if len(rawPassword) == 0 {
		return nil, ErrEmptyPassword
	}

	password := Password{password: rawPassword}

	err := password.Validate()
	if err != nil {
		return nil, err
	}

	return &password, nil
}

func (p *Password) Validate() error {
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
