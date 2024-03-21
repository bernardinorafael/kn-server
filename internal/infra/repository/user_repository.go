package repository

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) contract.UserRepository {
	return &userRepository{DB}
}

func (ur *userRepository) Save(u entity.User) error {
	user := entity.User{
		ID:       uuid.New().String(),
		Name:     u.Name,
		Email:    u.Email,
		Document: u.Document,
		Password: u.Password,
	}

	err := ur.DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) GetByEmail(email string) (*entity.User, error) {
	var user = entity.User{}

	err := ur.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetByID(id string) (*entity.User, error) {
	var user = entity.User{ID: id}

	if err := ur.DB.First(&user).Error; err != nil {
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
