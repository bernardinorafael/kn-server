package contract

import "github.com/bernardinorafael/kn-server/internal/domain/entity"

type UserRepository interface {
	Create(u entity.User) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
}
