package contract

import (
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Login(mail, password string) (*user.User, error)
	Register(name, mail, password, document string) (*user.User, error)
	RecoverPassword(id int, data dto.UpdatePassword) error
}

type JWTService interface {
	CreateToken(id string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type ProductService interface {
	Create(data dto.CreateProduct) error
	Delete(publicID string) error
	GetByPublicID(publicID string) (*product.Product, error)
	GetBySlug(slugInput string) (*product.Product, error)
	GetAll() ([]product.Product, error)
}
