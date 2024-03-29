package service

import (
	"errors"
	"fmt"
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

func (js *jwtService) CreateToken(id string) (string, *dto.Claims, error) {
	duration := js.s.cfg.AccessTokenDuration
	secret := []byte(js.s.cfg.JWTSecret)

	claims := &dto.Claims{
		UserID:    id,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		js.s.log.Error("error to encrypt token", err)
		return "", claims, ErrEncryptToken
	}

	return token, claims, nil
}

func (js *jwtService) Decode(token string) (string, error) {
	if strings.TrimSpace(token) == "" {
		return "", ErrCouldNotParseJWT
	}

	secret := js.s.cfg.JWTSecret
	var id string

	keyFn := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidCredentials
		}
		return []byte(secret), nil
	}

	jwtToken, err := jwt.Parse(token, keyFn)
	if err != nil {
		return "", err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		id = fmt.Sprint(claims["id"])
	}

	return id, nil
}

func (js *jwtService) ValidateToken(token string) (*dto.Claims, error) {
	secret := js.s.cfg.JWTSecret

	if strings.TrimSpace(token) == "" {
		return nil, ErrInvalidCredentials
	}

	keyFn := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidCredentials
		}
		return []byte(secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &dto.Claims{}, keyFn)
	if err != nil {
		var v *jwt.ValidationError
		ok := errors.As(err, &v)
		if ok && strings.Contains(v.Inner.Error(), ErrExpiredToken.Error()) {
			js.s.log.Error(ErrExpiredToken.Error())
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*dto.Claims)
	if !ok {
		js.s.log.Error("could not parse jwt token", err)
		return nil, ErrCouldNotParseJWT
	}

	return claims, nil
}
