package contract

import (
	"context"

	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type AccountService interface {
	CreateAccount(ctx context.Context, i dto.CreateAccount) error
	GetByID(ctx context.Context, id string) (*entity.Account, error)
	UpdateAccount(ctx context.Context, i dto.UpdateAccount, id string) error
	DeleteAccount(ctx context.Context, id string) error
	GetAll(ctx context.Context) (*[]entity.Account, error)
}

type AuthService interface {
	Register(ctx context.Context, i dto.Register) (id string, err error)
	Login(ctx context.Context, i dto.Login) (id string, err error)
}

type JWTService interface {
	CreateToken(ctx context.Context, id string) (string, *dto.Claims, error)
	ValidateToken(ctx context.Context, token string) (*dto.Claims, error)
}
