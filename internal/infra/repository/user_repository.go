package repository

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) contract.UserRepository {
	return &userRepository{DB}
}

func (repository *userRepository) Create(u entity.User) error {
	user := entity.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}

	err := repository.DB.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *userRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := repository.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *userRepository) GetByID(id string) (*entity.User, error) {
	var user = entity.User{ID: id}

	if err := repository.DB.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetAll() (*[]entity.User, error) {
	var users []entity.User

	if err := ur.DB.First(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &users, nil
}

func (ur *userRepository) GetPassword(id string) (string, error) {
	var user = entity.User{ID: id}

	err := ur.DB.Select("password").Find(&user).Error
	if err != nil {
		return "", err
	}

	return user.Password, nil
}
