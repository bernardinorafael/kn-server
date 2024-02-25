package interfaces

import (
	"github.com/bernardinorafael/gozinho/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Save(u *entity.User) error
	GetByEmail(email string) (*entity.User, error)
	GetByID(id uuid.UUID) (*entity.User, error)
	Update(u *entity.User) error
	Delete(id uuid.UUID) error
	GetAll() ([]*entity.User, error)
	UpdatePassword(password string, id uuid.UUID) error
}
