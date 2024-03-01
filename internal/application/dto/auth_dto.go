package dto

import (
	"errors"
	"time"
)

type TokenPayloadInput struct {
	ID string
}

type TokenPayload struct {
	UserID    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func (t TokenPayload) Valid() error {
	if time.Now().After(t.ExpiresAt) {
		return errors.New("the provided access token has expired")
	}
	return nil
}
