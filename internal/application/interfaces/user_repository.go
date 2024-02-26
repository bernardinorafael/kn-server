package interfaces

import (
	"github.com/bernardinorafael/gozinho/internal/domain/entity"
)

type UserRepository interface {
	Save(u *entity.User) error
	GetByEmail(email string) (*entity.User, error)
	GetByID(id string) (*entity.User, error)
	Update(u *entity.User) error
	Delete(id string) error
	GetAll() ([]entity.User, error)
	UpdatePassword(password string, id string) error
}
