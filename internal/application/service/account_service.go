package service

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/gorm"
)

type accountService struct {
	s *service
}

func newAccountService(service *service) contract.AccountService {
	return &accountService{s: service}
}

func (as *accountService) GetByID(id string) (*entity.Account, error) {
	as.s.log.Info("Process started")
	defer as.s.log.Info("Process finished")

	user, err := as.s.accountRepo.GetByID(id)
	if err != nil {
		as.s.log.Error("error to find user: %s", err.Error())
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (as *accountService) UpdateAccount(i dto.UpdateAccount, id string) error {
	as.s.log.Info("Process started")
	defer as.s.log.Info("Process finished")

	if i.Email != "" {
		_, err := as.s.accountRepo.GetByEmail(i.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEmailNotFound
		} else if err == nil {
			return ErrEmailAlreadyTaken
		}
	}

	account := entity.Account{
		Name:  i.Name,
		Email: i.Email,
	}

	err := as.s.accountRepo.Update(&account, id)
	if err != nil {
		as.s.log.Error("error update user: %s", err.Error())
		return ErrUpdateUser
	}

	return nil
}

func (as *accountService) DeleteAccount(id string) error {
	as.s.log.Info("Process started")
	defer as.s.log.Info("Process finished")

	_, err := as.s.accountRepo.GetByID(id)
	if err != nil {
		as.s.log.Error("error find user by ID: %s", err.Error())
		return ErrUserNotFound
	}

	err = as.s.accountRepo.Delete(id)
	if err != nil {
		as.s.log.Error("error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	return nil
}

func (as *accountService) GetAll() (*[]entity.Account, error) {
	as.s.log.Info("Process started")
	defer as.s.log.Info("Process finished")

	accounts, err := as.s.accountRepo.GetAll()
	if err != nil {
		if len(*accounts) == 0 {
			return nil, ErrEmptyResourceError
		}
		as.s.log.Error("error find users: %s", err.Error())
		return nil, ErrGetManyUsers
	}

	return accounts, nil
}
