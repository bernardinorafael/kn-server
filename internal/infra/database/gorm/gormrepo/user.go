package gormrepo

import (
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(DB *gorm.DB) contract.UserRepository {
	return &userRepo{db: DB}
}

func (r *userRepo) Create(u user.User) (*gormodel.User, error) {
	user := gormodel.User{
		PublicID:  u.PublicID(),
		Name:      u.Name(),
		Email:     string(u.Email()),
		Password:  string(u.Password()),
		Document:  string(u.Document()),
		Phone:     string(u.Phone()),
		UpdatedAt: time.Now(),
	}

	err := r.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByPublicID(publicID string) (*gormodel.User, error) {
	var u gormodel.User

	err := r.db.Where("public_id = ?", publicID).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByID(id int) (*gormodel.User, error) {
	var user gormodel.User

	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(email string) (*gormodel.User, error) {
	var u gormodel.User

	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(u user.User) (*gormodel.User, error) {
	var user gormodel.User

	updated := gormodel.User{
		PublicID:  u.PublicID(),
		Name:      u.Name(),
		Email:     string(u.Email()),
		Password:  string(u.Password()),
		Document:  string(u.Document()),
		Phone:     string(u.Phone()),
		UpdatedAt: time.Now(),
	}

	err := r.db.
		Model(&user).
		Where("public_id = ?", u.PublicID()).
		Updates(updated).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
