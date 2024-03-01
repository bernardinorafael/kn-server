package contract

import (
	"context"

	"github.com/bernardinorafael/kn-server/internal/application/dto"
)

type AuthService interface {
	CreateAccessToken(ctx context.Context, i dto.TokenPayloadInput) (string, *dto.TokenPayload, error)
	ValidateToken(ctx context.Context, token string) (*dto.TokenPayload, error)
}
