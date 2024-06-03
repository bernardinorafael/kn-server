package contract

import (
	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
)

type UserRepository interface {
	Create(u user.User) (*user.User, error)
	FindByID(id uint) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}

type ProductRepository interface {
	Create(product product.Product) (*product.Product, error)
	FindByID(id uint) (*product.Product, error)
	FindBySlug(name string) (*product.Product, error)
}
