package repository

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) contract.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(u entity.User) error {
	user := entity.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}

	err := r.db.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) FindByID(id string) (*entity.User, error) {
	var user entity.User

	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
