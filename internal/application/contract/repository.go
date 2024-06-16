package contract

import (
	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
)

type UserRepository interface {
	Create(u user.User) (*user.User, error)
	FindByID(id int) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
	Update(u user.User) (*user.User, error)
}

type ProductRepository interface {
	Create(product product.Product) (*product.Product, error)
	GetByID(id int) (*product.Product, error)
	GetByPublicID(publicID string) (*product.Product, error)
	GetBySlug(name string) (*product.Product, error)
	GetAll() ([]product.Product, error)
	Delete(publicID string) error
}
