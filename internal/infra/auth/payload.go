package auth

import (
	"errors"
	"time"
)

type PayloadInput struct {
	UserID    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func newPayload(userID string, duration time.Duration) *PayloadInput {
	return &PayloadInput{
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
}

func (c *PayloadInput) Valid() error {
	if time.Now().After(c.ExpiresAt) {
		return errors.New("token expired")
	}
	return nil
}
