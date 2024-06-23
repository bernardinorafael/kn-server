package contract

import (
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
)

type AuthService interface {
	Login(mail, password string) (*user.User, error)
	Register(name, mail, password, document string) (*user.User, error)
	RecoverPassword(publicID string, data dto.UpdatePassword) error
}

type ProductService interface {
	Create(data dto.CreateProduct) error
	Delete(publicID string) error
	GetByPublicID(publicID string) (*product.Product, error)
	GetBySlug(slugInput string) (*product.Product, error)
	GetAll() ([]product.Product, error)
	UpdatePrice(publicID string, price float64) error
}

type UserService interface {
	GetUser(publicID string) (*user.User, error)
}
