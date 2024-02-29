package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bernardinorafael/gozinho/internal/application/contract"
	utillog "github.com/bernardinorafael/gozinho/util/log"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/chacha20poly1305"
)

type jwtAuthentication struct {
	jwtSecretKey string
	log          utillog.Logger
}

func newJwtAuthentication(jwtSecretKey string, log utillog.Logger) (contract.Authentication, error) {
	if len(jwtSecretKey) != chacha20poly1305.KeySize {
		return nil, errors.New("invalid token")
	}
	return &jwtAuthentication{
		jwtSecretKey, log,
	}, nil
}

func (a jwtAuthentication) GenerateToken(ctx context.Context, payload PayloadInput) (string, *PayloadInput, error) {
	return a.generateToken(ctx, payload, jwtTokenExpiresAt)
}

func (a jwtAuthentication) VerifyToken(ctx context.Context, token string) (*PayloadInput, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token invalido")
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token invalido")
		}
		return []byte(a.jwtSecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &PayloadInput{}, keyFunc)
	if err != nil {
		var val *jwt.ValidationError
		ok := errors.As(err, &val)
		if ok && strings.Contains(val.Inner.Error(), "expired token") {
			a.log.Errorf(ctx, "expired token: %v", err)
			return nil, err
		}
		return nil, err
	}

	var ok bool
	payload, ok := jwtToken.Claims.(*PayloadInput)
	if !ok {
		a.log.Errorf(ctx, "could not parse jwt token: %v", err)
		return nil, errors.New("could not parse jwt token")
	}
	return payload, nil
}

func (a jwtAuthentication) generateToken(ctx context.Context, payload PayloadInput, d time.Duration) (string, *PayloadInput, error) {
	key := []byte(a.jwtSecretKey)
	claims := newPayload(payload.UserID, d)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString(key)
	if err != nil {
		a.log.Errorf(ctx, "error to encrypt token: %v", err.Error())
		return "", claims, err
	}

	return token, claims, nil
}
