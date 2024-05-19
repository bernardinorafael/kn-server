package service

import (
	"log/slog"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
)

type jwtService struct {
	log *slog.Logger
	env *config.EnvFile
}

func NewJWTService(log *slog.Logger, env *config.EnvFile) contract.JWTService {
	return &jwtService{log, env}
}

func (j *jwtService) CreateToken() (string, error) {
	panic("unimplemented")
}
