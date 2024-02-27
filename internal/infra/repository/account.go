package repository

import (
	"github.com/bernardinorafael/gozinho/internal/application/contract"
	"github.com/bernardinorafael/gozinho/internal/domain/entity"
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
		ID:         uuid.New().String(),
		Name:       u.Name,
		Username:   u.Username,
		Email:      u.Email,
		PersonalID: u.PersonalID,
		Password:   u.Password,
	}

	err := r.DB.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetByEmail(email string) (*entity.Account, error) {
	var user = entity.Account{}

	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *accountRepository) GetByID(id string) (*entity.Account, error) {
	var user = entity.Account{ID: id}

	if err := r.DB.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *accountRepository) Update(u *entity.Account) error {
	var user = entity.Account{ID: u.ID}

	if err := r.DB.First(&user).Error; err != nil {
		return err
	}

	if u.Name != "" {
		user.Name = u.Name
	}
	if u.Username != "" {
		user.Username = u.Username
	}
	if u.Email != "" {
		user.Email = u.Email
	}

	if err := r.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetAll() ([]entity.Account, error) {
	var users []entity.Account

	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *accountRepository) Delete(id string) error {
	var user = entity.Account{ID: id}

	if err := r.DB.First(&user).Error; err != nil {
		return err
	}

	if err := r.DB.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) UpdatePassword(password string, id string) error {
	var user = entity.Account{ID: id}

	if err := r.DB.First(&user).Error; err != nil {
		return err
	}

	err := r.DB.Model(&user).Update("password", password).Error
	if err != nil {
		return err
	}

	return nil
}
