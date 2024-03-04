package contract

import (
	"context"

	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type AccountService interface {
	GetByID(id string) (*entity.Account, error)
	UpdateAccount(i dto.UpdateAccount, id string) error
	DeleteAccount(id string) error
	GetAll() (*[]entity.Account, error)
}

type AuthService interface {
	Register(ctx context.Context, i dto.Register) (*entity.Account, error)
	Login(ctx context.Context, i dto.Login) (*entity.Account, error)
}

type JWTService interface {
	CreateToken(id string) (string, *dto.Claims, error)
	ValidateToken(token string) (*dto.Claims, error)
}
