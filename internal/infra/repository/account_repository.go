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
	account := entity.User{
		ID:       uuid.New().String(),
		Name:     u.Name,
		Email:    u.Email,
		Document: u.Document,
		Password: u.Password,
	}

	err := r.DB.Create(&account).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetByEmail(email string) (*entity.User, error) {
	var account = entity.User{}

	err := r.DB.Where("email = ?", email).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *accountRepository) GetByID(id string) (*entity.User, error) {
	var account = entity.User{ID: id}

	if err := r.DB.First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *accountRepository) Update(account *entity.User, id string) error {
	err := r.DB.Debug().
		Model(account).
		Where("id = ?", id).
		Updates(account).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetAll() (*[]entity.User, error) {
	var accounts []entity.User

	if err := r.DB.First(&accounts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &accounts, nil
}

func (r *accountRepository) Delete(id string) error {
	var account = entity.User{ID: id}

	if err := r.DB.First(&account).Error; err != nil {
		return err
	}

	if err := r.DB.Delete(&account).Error; err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) UpdatePassword(password string, id string) error {
	var account = entity.User{ID: id}

	if err := r.DB.First(&account).Error; err != nil {
		return err
	}

	err := r.DB.Model(&account).Update("password", password).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetPassword(id string) (string, error) {
	var account = entity.User{ID: id}

	err := r.DB.Select("password").Find(&account).Error
	if err != nil {
		return "", err
	}

	return account.Password, nil
}
