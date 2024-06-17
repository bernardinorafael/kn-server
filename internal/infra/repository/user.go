package repository

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(DB *gorm.DB) contract.UserRepository {
	return &userRepo{db: DB}
}

func (r *userRepo) Create(u user.User) (*user.User, error) {
	user := &user.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		PublicID: u.PublicID,
		Document: u.Document,
	}

	err := r.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) GetByPublicID(publicID string) (*user.User, error) {
	var user user.User

	err := r.db.Where("public_id = ?", publicID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByID(id int) (*user.User, error) {
	var user user.User

	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(email string) (*user.User, error) {
	var user user.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(u user.User) (*user.User, error) {
	user := user.User{}

	updated := map[string]interface{}{
		"Name":     u.Name,
		"Password": u.Password,
	}

	err := r.db.
		Model(&user).
		Where("public_id = ?", u.PublicID).
		Updates(updated).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}
