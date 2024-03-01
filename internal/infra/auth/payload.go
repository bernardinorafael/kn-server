package auth

import (
	"errors"
	"time"

	"github.com/bernardinorafael/kn-server/internal/application/dto"
)

type TokenPayload struct {
	UserID    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func newPayload(tk dto.TokenPayloadInput, d time.Duration) *TokenPayload {
	return &TokenPayload{
		UserID:    tk.ID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(d),
	}
}

func (p *TokenPayload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return errors.New("the provided access token has expired")
	}
	return nil
}
