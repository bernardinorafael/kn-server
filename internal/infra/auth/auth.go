package auth

import (
	"context"
	"errors"
	"time"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	utillog "github.com/bernardinorafael/kn-server/util/log"
)

type Key string

const TokenKey Key = "Authorization"
const UserIDKey Key = "user-id"

var jwtTokenExpiresAt time.Duration

var (
	expiredTokenError     = errors.New("the provided access token has expired")
	invalidTokenError     = errors.New("the provided access token is invalid")
	couldNotParseJwtError = errors.New("failed to parse the provided jwt token")
	encryptTokenError     = errors.New("failed to encrypt the provided token")
)

type Authentication interface {
	CreateAccessToken(ctx context.Context, i dto.TokenPayloadInput) (string, *TokenPayload, error)
	ValidateToken(ctx context.Context, token string) (*TokenPayload, error)
}

func New(cfg *config.EnvFile, l utillog.Logger) (Authentication, error) {
	jwtTokenExpiresAt = cfg.AccessTokenDuration
	return newJwtAuthentication(cfg.JWTSecret, l)
}
