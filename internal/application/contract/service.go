package contract

import (
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Login(email, password string) (*entity.User, error)
	Register(name, email, password string) (*entity.User, error)
}

type JWTService interface {
	CreateToken(id uint) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type ProductService interface {
	Create(data dto.CreateProduct) error
}
