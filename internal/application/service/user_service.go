package service

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/gorm"
)

type userService struct {
	s *service
}

func newUserService(service *service) contract.UserService {
	return &userService{s: service}
}

func (as *userService) GetByID(id string) (*entity.User, error) {
	as.s.log.Info("Process started")
	defer as.s.log.Info("Process finished")

	user, err := as.s.userRepo.GetByID(id)
	if err != nil {
		as.s.log.Error("error to find user", err)
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (as *userService) UpdateUser(i dto.UpdateUser, id string) error {
	as.s.log.Info("Process started")
	defer as.s.log.Info("Process finished")

	if i.Email != "" {
		_, err := as.s.userRepo.GetByEmail(i.Email)
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

	err := as.s.userRepo.Update(&account, id)
	if err != nil {
		as.s.log.Error("error update user", err)
		return ErrUpdateUser
	}

	return nil
}

func (as *userService) DeleteUser(id string) error {
	as.s.log.Info("Process started")
	defer as.s.log.Info("Process finished")

	_, err := as.s.userRepo.GetByID(id)
	if err != nil {
		as.s.log.Error("error find user by ID", err)
		return ErrUserNotFound
	}

	err = as.s.userRepo.Delete(id)
	if err != nil {
		as.s.log.Error("error deleting user", err)
		return ErrDeleteUser
	}

	return nil
}

func (as *userService) GetAll() (*[]entity.User, error) {
	as.s.log.Info("Process started")
	defer as.s.log.Info("Process finished")

	accounts, err := as.s.userRepo.GetAll()
	if err != nil {
		if len(*accounts) == 0 {
			return nil, ErrEmptyResourceError
		}
		as.s.log.Error("error find users", err)
		return nil, ErrGetManyUsers
	}

	return accounts, nil
}
