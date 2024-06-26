package contract

import (
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
)

type UserRepository interface {
	Create(usr user.User) (*user.User, error)
	GetByID(id int) (*user.User, error)
	GetByPublicID(publicID string) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Update(usr user.User) (*user.User, error)
}

type ProductRepository interface {
	Create(prod product.Product) (*product.Product, error)
	GetByID(id int) (*product.Product, error)
	GetByPublicID(publicID string) (*product.Product, error)
	GetBySlug(name string) (*product.Product, error)
	GetAll() ([]product.Product, error)
	Delete(publicID string) error
	Update(prod product.Product) error
}
