package email

import (
	"errors"
	"regexp"
	"strings"
)

const (
	maxAddressLength    = 254
	minAddressLength    = 3
	minDomainPartLength = 3
)

var (
	ErrEmptyEmailAddress       = errors.New("email cannot be empty")
	ErrInvalidEmailFormat      = errors.New("invalid email address")
	ErrEmailTooLong            = errors.New("email address is too long")
	ErrEmailContainsWhiteSpace = errors.New("email cannot contain whitespace")
	ErrEmailTooShort           = errors.New("local part of the email must have at least 3 characters")
	ErrDomainPartTooShort      = errors.New("domain part email must have at least 3 characters")
)

type Email string

type Address struct {
	local  string
	domain string
}

func New(address string) (*Address, error) {
	if len(address) == 0 {
		return nil, errors.New("email is a required field")
	}

	splitted := strings.Split(address, "@")

	if len(splitted) != 2 {
		return nil, ErrInvalidEmailFormat
	}

	addr := Address{local: splitted[0], domain: splitted[1]}

	err := addr.validate()
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func (a *Address) validate() error {
	email := a.Email()

	if strings.TrimSpace(string(email)) == "" {
		return ErrEmptyEmailAddress
	}

	if strings.ContainsAny(string(email), "\n\t") {
		return ErrEmailContainsWhiteSpace
	}

	if len(email) > maxAddressLength {
		return ErrEmailTooLong
	}

	if len(a.local) < minAddressLength {
		return ErrEmailTooShort
	}

	localPartPattern := `^[a-zA-Z0-9][a-zA-Z0-9._-]*[a-zA-Z0-9]$`
	matched, _ := regexp.MatchString(localPartPattern, a.local)
	if !matched {
		return ErrInvalidEmailFormat
	}

	if strings.Contains(a.local, "..") {
		return ErrInvalidEmailFormat
	}

	domainPartPattern := `^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`
	matched, _ = regexp.MatchString(domainPartPattern, a.domain)
	if !matched {
		return ErrInvalidEmailFormat
	}

	splittedDomainPart := strings.Split(a.domain, ".")
	if len(splittedDomainPart[0]) < minDomainPartLength {
		return ErrDomainPartTooShort
	}

	if len(splittedDomainPart) < 2 {
		return ErrInvalidEmailFormat
	}

	return nil
}

func (a *Address) GetLocalPart() Email {
	return Email(a.local)
}

func (a *Address) GetDomainPart() Email {
	return Email(a.domain)
}

func (a *Address) Email() Email {
	return Email(a.local + "@" + a.domain)
}
