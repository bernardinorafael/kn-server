package repository

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type accountRepository struct {
	DB *gorm.DB
}

func NewAccountRepository(DB *gorm.DB) contract.UserRepository {
	return &accountRepository{DB}
}

func (r *accountRepository) Save(u entity.User) error {
	user := entity.User{
		ID:       uuid.New().String(),
		Name:     u.Name,
		Email:    u.Email,
		Document: u.Document,
		Password: u.Password,
	}

	err := r.DB.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetByEmail(email string) (*entity.User, error) {
	var user = entity.User{}

	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *accountRepository) GetByID(id string) (*entity.User, error) {
	var user = entity.User{ID: id}

	if err := r.DB.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *accountRepository) Update(user *entity.User, id string) error {
	err := r.DB.
		Model(&user).
		Where("id = ?", id).
		Updates(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetAll() (*[]entity.User, error) {
	var users []entity.User

	if err := r.DB.First(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &users, nil
}

func (r *accountRepository) Delete(id string) error {
	var err error
	var user = entity.User{ID: id}

	if err = r.DB.First(&user).Error; err != nil {
		return err
	}

	if err = r.DB.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) UpdatePassword(password string, id string) error {
	var user = entity.User{ID: id}

	if err := r.DB.First(&user).Error; err != nil {
		return err
	}

	err := r.DB.Model(&user).Update("password", password).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetPassword(id string) (string, error) {
	var user = entity.User{ID: id}

	err := r.DB.Select("password").Find(&user).Error
	if err != nil {
		return "", err
	}

	return user.Password, nil
}
