package repository

import (
	"github.com/bernardinorafael/gozinho/internal/application/contract"
	"github.com/bernardinorafael/gozinho/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) contract.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Save(u *entity.User) error {
	user := entity.User{
		ID:         uuid.New().String(),
		Name:       u.Name,
		Username:   u.Username,
		Email:      u.Email,
		PersonalID: u.PersonalID,
		Password:   u.Password,
	}

	err := r.db.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*entity.User, error) {
	var user = entity.User{}

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *userRepository) GetByID(id string) (*entity.User, error) {
	var user = entity.User{ID: id}

	if err := r.db.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(u *entity.User) error {
	var user = entity.User{ID: u.ID}

	if err := r.db.First(&user).Error; err != nil {
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

	if err := r.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetAll() ([]entity.User, error) {
	var users []entity.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Delete(id string) error {
	var user = entity.User{ID: id}

	if err := r.db.First(&user).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdatePassword(password string, id string) error {
	var user = entity.User{ID: id}

	if err := r.db.First(&user).Error; err != nil {
		return err
	}

	err := r.db.Model(&user).Update("password", password).Error
	if err != nil {
		return err
	}

	return nil
}
