package auth

import (
	"errors"
	"time"
)

var (
	errInvalidSecret = errors.New("invalid secret key")
	errEncryptSecret = errors.New("failed to encrypt secret key")
	errInvalidToken  = errors.New("invalid token")
	errUnauthorized  = errors.New("unauthorized")
)

type TokenAuthInterface interface {
	CreateAccessToken(publicID string, duration time.Duration) (token string, payload *TokenPayload, err error)
	VerifyToken(t string) (*TokenPayload, error)
}

type TokenPayload struct {
	PublicID  string    `json:"public_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func newPayload(publicID string, duration time.Duration) *TokenPayload {
	return &TokenPayload{
		PublicID:  publicID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
}

func (j *TokenPayload) Valid() error {
	if time.Now().After(j.ExpiredAt) {
		return errors.New("expired token")
	}
	return nil
}
