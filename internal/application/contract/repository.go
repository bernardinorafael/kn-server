package contract

import "github.com/bernardinorafael/kn-server/internal/domain/entity"

type UserRepository interface {
	Create(u entity.User) error
	FindByID(id string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}
