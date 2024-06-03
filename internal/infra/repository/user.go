package repository

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
	"gorm.io/gorm"
)

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepo(DB *gorm.DB) contract.UserRepository {
	return &userRepo{DB: DB}
}

func (r *userRepo) Create(u user.User) (*user.User, error) {
	user := &user.User{
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

func (r *userRepo) FindByID(id uint) (*user.User, error) {
	var user user.User

	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindByEmail(email string) (*user.User, error) {
	var user user.User

	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
