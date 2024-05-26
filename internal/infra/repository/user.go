package repository

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/gorm"
)

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepo(DB *gorm.DB) contract.UserRepository {
	return &userRepo{DB: DB}
}

func (r *userRepo) Create(u entity.User) (*entity.User, error) {
	user := &entity.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}

	err := r.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) FindByID(id string) (*entity.User, error) {
	var user entity.User

	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
