package money

import (
	"errors"
)

const (
	maxCentsAllowed = 99999999 // R$ 999.999,99
	minCentsAllowed = 1
)

type Money int

type money struct {
	value int
}

func New(cents int) (money, error) {
	m := money{cents}
	if err := m.validate(); err != nil {
		return money{}, err
	}
	return m, nil
}

func (m *money) validate() error {
	if m.value < minCentsAllowed {
		return errors.New("monetary amounts cannot be negative")
	}

	if m.value > maxCentsAllowed {
		return errors.New("maximum amount value has been exceeded")
	}

	return nil
}

func (m *money) Cents() Money { return Money(m.value) }
