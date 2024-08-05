package cpf

import (
	"errors"
	"fmt"
	"strconv"
)

type CPF string

type cpf struct {
	value string
}

func New(document string) (*cpf, error) {
	if len(document) == 0 {
		return nil, errors.New("document is a required field")
	}

	doc := cpf{value: document}

	if err := doc.validate(); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (c *cpf) validate() error {
	if len(c.value) != 11 {
		return fmt.Errorf("cpf must have 11 characters")
	}

	s1 := sum(c.value[:9], []int{10, 9, 8, 7, 6, 5, 4, 3, 2})
	d1 := applyValidationRule(s1)

	if applyValidationRule(s1) >= 10 {
		d1 = 0
	}

	s2 := sum(fmt.Sprintf("%v%v", c.value[:9], d1), []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2})
	d2 := applyValidationRule(s2)

	if applyValidationRule(s2) >= 10 {
		d2 = 0
	}

	validated := fmt.Sprintf("%v%v%v", c.value[:9], d1, d2)

	if validated != c.value {
		return fmt.Errorf("invalid cpf")
	}
	return nil
}

func sum(cpf string, table []int) int {
	sum := 0
	for i, v := range table {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * v
	}
	return sum
}

func applyValidationRule(sum int) int {
	return 11 - (sum % 11)
}

func (c *cpf) CPF() CPF {
	return CPF(c.value)
}
