package auth

import (
	"time"

	"github.com/bernardinorafael/gozinho/config"
	"github.com/bernardinorafael/gozinho/internal/application/contract"
	utillog "github.com/bernardinorafael/gozinho/util/log"
)

var jwtTokenExpiresAt time.Duration

func NewAuthenticationToken(cfg *config.EnvFile, l utillog.Logger) (contract.Authentication, error) {
	jwtTokenExpiresAt = cfg.AccessTokenDuration
	return newJwtAuthentication(cfg.JWTSecret, l)
}
