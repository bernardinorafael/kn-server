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

func NewAccountRepository(DB *gorm.DB) contract.AccountRepository {
	return &accountRepository{DB}
}

func (r *accountRepository) Save(u *entity.Account) error {
	user := entity.Account{
		ID:       uuid.New().String(),
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Document: u.Document,
		Password: u.Password,
	}

	if err := r.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) CheckUserExist(email, username, Document string) (bool, error) {
	var user entity.Account

	err := r.DB.Where("email = ? OR username = ? OR document = ?",
		email, username, Document).First(&user).Error

	if err == nil {
		return true, nil
	}
	return false, err
}

func (r *accountRepository) GetByEmail(email string) (*entity.Account, error) {
	var account = entity.Account{}

	err := r.DB.Where("email = ?", email).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, err
}

func (r *accountRepository) GetByID(id string) (*entity.Account, error) {
	var account = entity.Account{ID: id}

	if err := r.DB.First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *accountRepository) Update(u *entity.Account) error {
	var account = entity.Account{ID: u.ID}

	if err := r.DB.First(&account).Error; err != nil {
		return err
	}

	if u.Name != "" {
		account.Name = u.Name
	}
	if u.Username != "" {
		account.Username = u.Username
	}
	if u.Email != "" {
		account.Email = u.Email
	}

	if err := r.DB.Save(&account).Error; err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetAll() ([]entity.Account, error) {
	var accounts []entity.Account

	if err := r.DB.First(&accounts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return accounts, nil
}

func (r *accountRepository) Delete(id string) error {
	var account = entity.Account{ID: id}

	if err := r.DB.First(&account).Error; err != nil {
		return err
	}

	if err := r.DB.Delete(&account).Error; err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) UpdatePassword(password string, id string) error {
	var account = entity.Account{ID: id}

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
	var account = entity.Account{ID: id}

	err := r.DB.First(&account).Error
	if err != nil {
		return "", err
	}
	p := account.Password

	return p, nil
}
