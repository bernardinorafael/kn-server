package contract

import (
	"context"

	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type AccountService interface {
	GetByID(ctx context.Context, id string) (*entity.Account, error)
	UpdateAccount(ctx context.Context, i dto.UpdateAccount, id string) error
	DeleteAccount(ctx context.Context, id string) error
	GetAll(ctx context.Context) (*[]entity.Account, error)
}

type AuthService interface {
	Register(ctx context.Context, i dto.Register) (*entity.Account, error)
	Login(ctx context.Context, i dto.Login) (*entity.Account, error)
}

type JWTService interface {
	CreateToken(ctx context.Context, id string) (string, *dto.Claims, error)
	ValidateToken(ctx context.Context, token string) (*dto.Claims, error)
}
