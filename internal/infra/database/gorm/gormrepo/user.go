package gormrepo

import (
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"gorm.io/gorm"
)

/*
* TODO: remove entity/model mapping logic from repositories and do it into service layer
 */

type userRepo struct{ db *gorm.DB }

func NewUserRepo(DB *gorm.DB) contract.GormUserRepository {
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

func (r userRepo) Create(u gormodel.User) error {
	err := r.db.Create(&u).Error
	if err != nil {
		return err
	}
	return nil
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

func (r userRepo) Update(u gormodel.User) (gormodel.User, error) {
	var user gormodel.User

	err := r.db.
		Where("public_id = ?", u.PublicID).
		First(&user).
		Error
	if err != nil {
		return user, err
	}

	user.Name = u.Name
	user.Email = u.Email
	user.Password = u.Password
	user.Phone = u.Phone
	user.Status = u.Status
	user.UpdatedAt = time.Now()

	err = r.db.Save(&user).Error
	if err != nil {
		return gormodel.User{}, err
	}

	return user, nil
}
