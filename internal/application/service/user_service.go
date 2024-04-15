package service

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type userService struct {
	s *service
}

func newUserService(service *service) contract.UserService {
	return &userService{s: service}
}

func (us *userService) GetByID(id string) (*entity.User, error) {
	us.s.log.Info("Process started")
	defer us.s.log.Info("Process finished")

	user, err := us.s.userRepo.GetByID(id)
	if err != nil {
		us.s.log.Error("error to find user", err)
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (us *userService) GetAll() (*[]entity.User, error) {
	us.s.log.Info("Process started")
	defer us.s.log.Info("Process finished")

	accounts, err := us.s.userRepo.GetAll()
	if err != nil {
		if len(*accounts) == 0 {
			return nil, ErrEmptyResourceError
		}
		us.s.log.Error("error find users", err)
		return nil, ErrGetManyUsers
	}

	return accounts, nil
}
