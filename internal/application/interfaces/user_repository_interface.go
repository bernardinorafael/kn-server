package interfaces

import (
	"github.com/bernardinorafael/gozinho/internal/domain/entities"
)

type UserRepository interface {
	Create(u *entities.User) error
}
