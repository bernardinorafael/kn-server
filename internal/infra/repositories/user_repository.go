package repositories

import (
	"github.com/bernardinorafael/gozinho/internal/application/interfaces"
	"github.com/bernardinorafael/gozinho/internal/domain/entities"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(u *entities.User) error {
	return nil
}
