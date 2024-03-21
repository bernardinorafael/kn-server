package contract

import (
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type UserService interface {
	GetByID(id string) (*entity.User, error)
	GetAll() (*[]entity.User, error)
}
