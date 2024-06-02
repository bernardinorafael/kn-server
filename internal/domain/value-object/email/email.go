package email

import (
	"errors"
	"regexp"
	"strings"
)

const (
	maxAddressLength   = 254
	minAddressLength   = 3
	minDomainPartLengh = 3
)

var (
	ErrEmptyEmailAddress       = errors.New("Email cannot be empty")
	ErrInvalidEmailFormat      = errors.New("Invalid email address")
	ErrEmailTooLong            = errors.New("Email address is too long")
	ErrEmailContainsWhiteSpace = errors.New("Email cannot contain whitespace")
	ErrInvalidEmailChar        = errors.New("Email has some invalid special characters")
	ErrEmailTooShort           = errors.New("Local part of the email must have at least 3 characters")
	ErrDomainPartTooShort      = errors.New("Domain part email must have at least 3 characters")
)

type Address struct {
	local  string
	domain string
}

func New(address string) (*Address, error) {
	splitted := strings.Split(address, "@")

	if len(splitted) != 2 {
		return nil, ErrInvalidEmailFormat
	}

	return &Address{
		local:  splitted[0],
		domain: splitted[1],
	}, nil
}

func (a *Address) Validate() error {
	email := a.GetString()

	if strings.TrimSpace(email) == "" {
		return ErrEmptyEmailAddress
	}

	if strings.ContainsAny(email, "\n\t") {
		return ErrEmailContainsWhiteSpace
	}

	if len(email) > maxAddressLength {
		return ErrEmailTooLong
	}

	if len(a.local) < minAddressLength {
		return ErrEmailTooShort
	}

	localPartRegexp := `^[a-zA-Z0-9][a-zA-Z0-9._-]*[a-zA-Z0-9]$`
	matched, _ := regexp.MatchString(localPartRegexp, a.local)
	if !matched {
		return ErrInvalidEmailChar
	}

	if strings.Contains(a.local, "..") {
		return ErrInvalidEmailFormat
	}

	domainPartRegexp := `^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`
	matched, _ = regexp.MatchString(domainPartRegexp, a.domain)
	if !matched {
		return ErrInvalidEmailChar
	}

	splittedDomainPart := strings.Split(a.domain, ".")
	if len(splittedDomainPart[0]) < minDomainPartLengh {
		return ErrDomainPartTooShort
	}

	return nil
}

func (a *Address) GetString() string {
	return a.local + "@" + a.domain
}

func (a *Address) GetLocalPart() string {
	return a.local
}

func (a *Address) GetDomainPart() string {
	return a.domain
}
