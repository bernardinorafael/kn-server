package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/golang-jwt/jwt"
)

type authService struct {
	service *service
}

func newAuthService(service *service) contract.AuthService {
	return &authService{service}
}

func (a *authService) CreateAccessToken(ctx context.Context, i dto.TokenPayloadInput) (string, *dto.TokenPayload, error) {
	secret := []byte(a.service.cfg.JWTSecret)

	payload := &dto.TokenPayload{
		UserID:    i.ID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(jwtTokenExpiresAt),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString(secret)
	if err != nil {
		a.service.l.Errorf(ctx, "error to encrypt token: %v", err.Error())
		return "", payload, encryptTokenError
	}
	return token, payload, nil
}

func (a *authService) ValidateToken(ctx context.Context, token string) (*dto.TokenPayload, error) {
	if strings.TrimSpace(token) == "" {
		return nil, invalidTokenError
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, invalidTokenError
		}
		return []byte(a.service.cfg.JWTSecret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &dto.TokenPayload{}, keyFunc)
	if err != nil {
		var v *jwt.ValidationError
		ok := errors.As(err, &v)
		if ok && strings.Contains(v.Inner.Error(), expiredTokenError.Error()) {
			a.service.l.Error(ctx, expiredTokenError.Error())
			return nil, expiredTokenError
		}
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*dto.TokenPayload)
	if !ok {
		a.service.l.Errorf(ctx, "could not parse jwt token: %v", err)
		return nil, couldNotParseJwtError
	}
	return payload, nil
}
