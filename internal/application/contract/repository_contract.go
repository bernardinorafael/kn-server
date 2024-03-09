package contract

import (
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type UserRepository interface {
	Save(u entity.User) error
	GetByEmail(email string) (*entity.User, error)
	GetByID(id string) (*entity.User, error)
	Update(account *entity.User, id string) error
	Delete(id string) error
	GetAll() (*[]entity.User, error)
	UpdatePassword(password string, id string) error
	GetPassword(id string) (string, error)
}

type ProductRepository interface {
	Save(p entity.Product) error
	GetAll() (*[]entity.Product, error)
	GetByName(p entity.Product) (*entity.Product, error)
}
