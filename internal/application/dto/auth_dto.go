package dto

import (
	"errors"
	"time"
)

type Register struct {
	Name     string `json:"name" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Document int    `json:"document" validate:"required,gte=0,len=10"`
	Password string `json:"password" validate:"required,min=8"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type Claims struct {
	UserID    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func (c Claims) Valid() error {
	if time.Now().After(c.ExpiresAt) {
		return errors.New("the provided access token has expired")
	}
	return nil
}
