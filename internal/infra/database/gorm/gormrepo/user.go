package gormrepo

import (
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(DB *gorm.DB) contract.UserRepository {
	return &userRepo{db: DB}
}

func (r *userRepo) Create(u user.User) (*model.User, error) {
	user := model.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		PublicID: u.PublicID,
		Document: u.Document,
		Phone:    u.Phone,
	}

	err := r.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByPublicID(publicID string) (*model.User, error) {
	var u model.User

	err := r.db.Where("public_id = ?", publicID).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByID(id int) (*model.User, error) {
	var user model.User

	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	var u model.User

	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(u user.User) (*model.User, error) {
	var user model.User

	updated := model.User{
		Name:      u.Name,
		Password:  u.Password,
		Document:  u.Document,
		Phone:     u.Phone,
		UpdatedAt: time.Now(),
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
