package contract

import (
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Login(email, password string) (*user.User, error)
	Register(name, email, password string) (*user.User, error)
	RecoverPassword(id int, data dto.UpdatePassword) error
}

type JWTService interface {
	CreateToken(id string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type ProductService interface {
	Create(data dto.CreateProduct) error
}
