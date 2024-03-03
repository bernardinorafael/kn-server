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

type jwtService struct {
	s *service
}

func newJWTService(service *service) contract.JWTService {
	return &jwtService{s: service}
}

func (js *jwtService) CreateToken(ctx context.Context, id string) (string, *dto.Claims, error) {
	duration := js.s.cfg.AccessTokenDuration
	secret := []byte(js.s.cfg.JWTSecret)

	claims := &dto.Claims{
		UserID:    id,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		js.s.log.Errorf(ctx, "error to encrypt token: %v", err.Error())
		return "", claims, ErrEncryptToken
	}

	return token, claims, nil
}

func (js *jwtService) ValidateToken(ctx context.Context, token string) (*dto.Claims, error) {
	if strings.TrimSpace(token) == "" {
		return nil, ErrInvalidCredentials
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidCredentials
		}
		return []byte(js.s.cfg.JWTSecret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &dto.Claims{}, keyFunc)
	if err != nil {
		var v *jwt.ValidationError
		ok := errors.As(err, &v)
		if ok && strings.Contains(v.Inner.Error(), ErrExpiredToken.Error()) {
			js.s.log.Error(ctx, ErrExpiredToken.Error())
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*dto.Claims)
	if !ok {
		js.s.log.Errorf(ctx, "could not parse jwt token: %v", err)
		return nil, ErrCouldNotParseJWT
	}

	return claims, nil
}
