package auth

import (
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/chacha20poly1305"
)

type jwtAuth struct {
	log       *slog.Logger
	secretKey string
}

func NewJWTAuth(log *slog.Logger, secretKey string) (TokenAuthInterface, error) {
	if len(secretKey) != chacha20poly1305.KeySize {
		return nil, errInvalidSecret
	}
	return &jwtAuth{log, secretKey}, nil
}

func (j *jwtAuth) CreateAccessToken(publicID string, duration time.Duration) (token string, payload *TokenPayload, err error) {
	payload = newPayload(publicID, duration)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err = t.SignedString([]byte(j.secretKey))
	if err != nil {
		j.log.Error("failed to encrypt jwt token")
		return "", nil, errEncryptSecret
	}
	return token, payload, nil
}

func (j *jwtAuth) VerifyToken(t string) (*TokenPayload, error) {
	if strings.TrimSpace(t) == "" {
		return nil, errInvalidToken
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errUnauthorized
		}
		return []byte(j.secretKey), nil
	}

	token, err := jwt.ParseWithClaims(t, &TokenPayload{}, keyFunc)
	if err != nil {
		j.log.Error("error validate jwt")
		return nil, errUnauthorized
	}

	payload, ok := token.Claims.(*TokenPayload)
	if !ok {
		j.log.Error("could not parse jwt token")
		return nil, errInvalidToken
	}
	return payload, nil
}
