package contract

import (
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type UserRepository interface {
	Save(u entity.User) error
	GetByEmail(email string) (*entity.User, error)
	GetByID(id string) (*entity.User, error)
	GetAll() (*[]entity.User, error)
	GetPassword(id string) (string, error)
}
