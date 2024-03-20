package service

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/gorm"
)

type userService struct {
	svc *service
}

func newUserService(service *service) contract.UserService {
	return &userService{svc: service}
}

func (us *userService) GetByID(id string) (*entity.User, error) {
	us.svc.log.Info("Process started")
	defer us.svc.log.Info("Process finished")

	user, err := us.svc.userRepo.GetByID(id)
	if err != nil {
		us.svc.log.Error("error to find user", err)
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (us *userService) UpdateUser(i dto.UpdateUser, id string) error {
	us.svc.log.Info("Process started")
	defer us.svc.log.Info("Process finished")

	if i.Email != "" {
		_, err := us.svc.userRepo.GetByEmail(i.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEmailNotFound
		} else if err == nil {
			return ErrEmailAlreadyTaken
		}
	}

	account := entity.User{
		Name:  i.Name,
		Email: i.Email,
	}

	err := us.svc.userRepo.Update(&account, id)
	if err != nil {
		us.svc.log.Error("error update user", err)
		return ErrUpdateUser
	}

	return nil
}

func (us *userService) DeleteUser(id string) error {
	us.svc.log.Info("Process started")
	defer us.svc.log.Info("Process finished")

	_, err := us.svc.userRepo.GetByID(id)
	if err != nil {
		us.svc.log.Error("error find user by ID", err)
		return ErrUserNotFound
	}

	err = us.svc.userRepo.Delete(id)
	if err != nil {
		us.svc.log.Error("error deleting user", err)
		return ErrDeleteUser
	}

	return nil
}

func (us *userService) GetAll() (*[]entity.User, error) {
	us.svc.log.Info("Process started")
	defer us.svc.log.Info("Process finished")

	accounts, err := us.svc.userRepo.GetAll()
	if err != nil {
		if len(*accounts) == 0 {
			return nil, ErrEmptyResourceError
		}
		us.svc.log.Error("error find users", err)
		return nil, ErrGetManyUsers
	}

	return accounts, nil
}
