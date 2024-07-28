package phone

import (
	"errors"
	"strings"
)

const (
	phoneLength = 11
)

var (
	areaCodes = []string{
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "21", "22", "24", "27", "28",
		"31", "32", "33", "34", "35", "37", "38", "41", "42", "43", "44", "45", "46", "47",
		"48", "49", "51", "53", "54", "55", "61", "62", "63", "64", "65", "66", "67", "68",
		"69", "71", "73", "74", "75", "77", "79", "81", "82", "83", "84", "85", "86", "87",
		"88", "89", "91", "92", "93", "94", "95", "96", "97", "98", "99",
	}
)

type Phone string

type phone struct {
	area  string
	local string
}

func New(value string) (*phone, error) {
	if value == "" {
		return nil, errors.New("phone cannot be empty")
	}

	area := value[:2]
	local := value[2:]

	p := &phone{area, local}

	if err := p.validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *phone) ToPhone() Phone {
	return Phone(p.area) + Phone(p.local)
}

func (p *phone) validate() error {
	if p.area == "" {
		return errors.New("invalid area code")
	}

	if p.local == "" {
		return errors.New("invalid local part")
	}

	if len(p.ToPhone()) != phoneLength {
		return errors.New("invalid phone format")
	}

	if err := p.verifyAreaCode(); err != nil {
		return err
	}

	if !strings.HasPrefix(p.local, "9") {
		return errors.New("invalid phone format")
	}

	return nil
}

func (p *phone) verifyAreaCode() error {
	for _, v := range areaCodes {
		if v == p.area {
			return nil
		}
	}
	return errors.New("invalid area code")
}
