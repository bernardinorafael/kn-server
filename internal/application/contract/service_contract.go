package contract

import (
	"context"

	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type UserService interface {
	GetByID(id string) (*entity.User, error)
	UpdateUser(i dto.UpdateAccount, id string) error
	DeleteUser(id string) error
	GetAll() (*[]entity.User, error)
}

type AuthService interface {
	Register(ctx context.Context, i dto.Register) (*entity.User, error)
	Login(ctx context.Context, i dto.Login) (*entity.User, error)
}

type JWTService interface {
	CreateToken(id string) (string, *dto.Claims, error)
	ValidateToken(token string) (*dto.Claims, error)
}
