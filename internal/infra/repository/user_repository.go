package repository

import (
	"github.com/bernardinorafael/gozinho/internal/application/interfaces"
	"github.com/bernardinorafael/gozinho/internal/domain/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Save(u *entity.User) error {
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*entity.User, error) {
	return nil, nil
}

func (r *UserRepository) GetByID(id string) (*entity.User, error) {
	return nil, nil
}

func (r *UserRepository) Update(user *entity.User) error {
	return nil
}

func (r *UserRepository) GetAll() ([]*entity.User, error) {
	return nil, nil
}

func (r *UserRepository) Delete(id string) error {
	return nil
}

func (r *UserRepository) UpdatePassword(password string, id string) error {
	return nil
}
