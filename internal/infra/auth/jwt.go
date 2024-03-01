package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/application/dto"
	utillog "github.com/bernardinorafael/kn-server/util/log"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/chacha20poly1305"
)

type jwtAuth struct {
	jwtSecret string
	l         utillog.Logger
}

func newJwtAuthentication(jwtSecret string, l utillog.Logger) (Authentication, error) {
	if len(jwtSecret) != chacha20poly1305.KeySize {
		return nil, invalidTokenError
	}
	return &jwtAuth{jwtSecret, l}, nil
}

func (a jwtAuth) CreateAccessToken(ctx context.Context, i dto.TokenPayloadInput) (string, *TokenPayload, error) {
	secret := []byte(a.jwtSecret)
	payload := newPayload(i, jwtTokenExpiresAt)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString(secret)
	if err != nil {
		a.l.Errorf(ctx, "error to encrypt token: %v", err.Error())
		return "", payload, encryptTokenError
	}
	return token, payload, nil
}

func (a jwtAuth) ValidateToken(ctx context.Context, token string) (*TokenPayload, error) {
	if strings.TrimSpace(token) == "" {
		return nil, invalidTokenError
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, invalidTokenError
		}
		return []byte(a.jwtSecret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &TokenPayload{}, keyFunc)
	if err != nil {
		var v *jwt.ValidationError
		ok := errors.As(err, &v)
		if ok && strings.Contains(v.Inner.Error(), expiredTokenError.Error()) {
			a.l.Error(ctx, expiredTokenError.Error())
			return nil, expiredTokenError
		}
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*TokenPayload)
	if !ok {
		a.l.Errorf(ctx, "could not parse jwt token: %v", err)
		return nil, couldNotParseJwtError
	}
	return payload, nil
}
