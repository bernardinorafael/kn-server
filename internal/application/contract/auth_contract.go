package contract

import (
	"context"

	"github.com/bernardinorafael/gozinho/internal/infra/auth"
)

type Authentication interface {
	GenerateToken(ctx context.Context, input auth.PayloadInput) (string, *auth.PayloadInput, error)
	VerifyToken(ctx context.Context, token string) (*auth.PayloadInput, error)
}
