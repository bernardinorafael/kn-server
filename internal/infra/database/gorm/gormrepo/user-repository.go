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

func (r *userRepo) Create(usr user.User) (*model.User, error) {
	u := &model.User{
		Name:     usr.Name,
		Email:    usr.Email,
		Password: usr.Password,
		PublicID: usr.PublicID,
		Document: usr.Document,
	}

	err := r.db.Create(&u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
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
	var u model.User

	err := r.db.Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	var u model.User

	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(usr user.User) (*model.User, error) {
	var u model.User

	updated := model.User{Name: usr.Name, Password: usr.Password, UpdatedAt: time.Now()}

	err := r.db.
		Model(&u).
		Where("public_id = ?", u.PublicID).
		Updates(updated).
		First(&u).
		Error

	if err != nil {
		return nil, err
	}
	return &u, nil
}
