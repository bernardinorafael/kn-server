package gormrepo

import (
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"gorm.io/gorm"
)

/*
* TODO: remove entity/model mapping logic from repositories and do it into service layer
 */

type userRepo struct{ db *gorm.DB }

func NewUserRepo(DB *gorm.DB) contract.UserRepository {
	return &userRepo{db: DB}
}

func (r userRepo) GetByPhone(phone string) (gormodel.User, error) {
	var user gormodel.User

	err := r.db.
		Where("phone = ?", phone).
		First(&user).
		Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r userRepo) Create(u user.User) (gormodel.User, error) {
	var user gormodel.User

	newUser := gormodel.User{
		PublicID:  u.PublicID(),
		Name:      u.Name(),
		Email:     string(u.Email()),
		Password:  string(u.Password()),
		Phone:     string(u.Phone()),
		UpdatedAt: time.Now(),
	}

	err := r.db.
		Create(&newUser).
		First(&user).
		Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r userRepo) GetByPublicID(publicID string) (gormodel.User, error) {
	var user gormodel.User

	err := r.db.
		Where("public_id = ?", publicID).
		First(&user).
		Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r userRepo) GetByEmail(email string) (gormodel.User, error) {
	var user gormodel.User

	err := r.db.
		Where("email = ?", email).
		First(&user).
		Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r userRepo) Update(u user.User) (gormodel.User, error) {
	var user gormodel.User

	err := r.db.
		Where("public_id = ?", u.PublicID()).
		First(&user).
		Error
	if err != nil {
		return user, err
	}

	user.Name = u.Name()
	user.Email = string(u.Email())
	user.Password = string(u.Password())
	user.Phone = string(u.Phone())
	user.UpdatedAt = time.Now()

	err = r.db.Save(&user).Error
	if err != nil {
		return gormodel.User{}, err
	}

	return user, nil
}
