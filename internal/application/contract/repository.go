package contract

import "github.com/bernardinorafael/kn-server/internal/domain/entity"

type UserRepository interface {
	Create(u entity.User) (*entity.User, error)
	FindByID(id uint) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}

type ProductRepository interface {
	Create(product entity.Product) (*entity.Product, error)
}
