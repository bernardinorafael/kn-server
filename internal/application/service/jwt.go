package service

import (
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

func (j *jwtService) CreateToken(id string) (string, error) {
	// TODO: add user id in token
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

// TODO: implement this one
func (j *jwtService) VerifyToken(token string) error {
	return nil
}
