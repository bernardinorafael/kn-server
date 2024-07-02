package contract

import (
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
)

type UserRepository interface {
	Create(usr user.User) (*model.User, error)
	GetByID(id int) (*model.User, error)
	GetByPublicID(publicID string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(usr user.User) (*model.User, error)
}

type ProductRepository interface {
	Create(prod product.Product) (*model.Product, error)
	GetByID(id int) (*model.Product, error)
	GetByPublicID(publicID string) (*model.Product, error)
	GetBySlug(name string) (*model.Product, error)
	GetAll() ([]model.Product, error)
	Delete(publicID string) error
	Update(prod product.Product) error
}
