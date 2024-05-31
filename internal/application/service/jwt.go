package service

import (
	"errors"
	"log/slog"
	"time"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/golang-jwt/jwt"
)

type jwtService struct {
	log *slog.Logger
	env *config.EnvFile
}

func NewJWTService(log *slog.Logger, env *config.EnvFile) contract.JWTService {
	return &jwtService{log, env}
}

func (j *jwtService) CreateToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         id,
		"issued_at":  time.Now(),
		"expires_at": time.Now().Add(j.env.AccessTokenDuration),
	})

	tokenString, err := token.SignedString([]byte(j.env.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (j *jwtService) ValidateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid jwt token")
		}
		return []byte(j.env.JWTSecret), nil
	})
}
