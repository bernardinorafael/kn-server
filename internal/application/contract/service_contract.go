package contract

import (
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type UserService interface {
	GetByID(id string) (*entity.User, error)
	UpdateUser(i dto.UpdateUser, id string) error
	DeleteUser(id string) error
	GetAll() (*[]entity.User, error)
}

type ProductService interface {
	SaveProduct(p entity.Product) error
	GetByID(id uint) (*entity.Product, error)
}

type AuthService interface {
	Register(i dto.Register) (*entity.User, error)
	Login(i dto.Login) (*entity.User, error)
}

type JWTService interface {
	CreateToken(id string) (string, *dto.Claims, error)
	ValidateToken(token string) (*dto.Claims, error)
}
