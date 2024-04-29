package repository

import (
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uuid.UUID) (*entity.User, error)
}
