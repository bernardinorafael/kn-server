package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bernardinorafael/kn-server/helper/crypto"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type accountService struct {
	s *service
}

func newAccountService(service *service) contract.AccountService {
	return &accountService{s: service}
}

func (as *accountService) CreateAccount(ctx context.Context, i dto.CreateAccount) error {
	as.s.log.Info(ctx, "Process started")
	defer as.s.log.Info(ctx, "Process finished")

	_, err := as.s.accountRepo.GetByEmail(i.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		as.s.log.Errorf(ctx, "error to get account by email: %s", err.Error())
		return ErrEmailNotFound
	} else if err == nil {
		as.s.log.Error(ctx, "email already taken")
		return ErrEmailAlreadyTaken
	}

	password, err := crypto.EncryptPassword(i.Password)
	if err != nil {
		as.s.log.Errorf(ctx, "error hashing password: %v", err.Error())
		return ErrHashPassword
	}

	user := entity.Account{
		ID:        uuid.New().String(),
		Name:      i.Name,
		Email:     i.Email,
		Document:  i.Document,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = as.s.accountRepo.Save(user)
	if err != nil {
		if strings.Contains(err.Error(), "uni_accounts_document") {
			as.s.log.Error(ctx, "document already exist")
			return ErrDocumentALreadyTaken
		}
		as.s.log.Error(ctx, "error creating account")
		return ErrCreateUser
	}

	return nil
}

func (as *accountService) GetByID(ctx context.Context, id string) (*entity.Account, error) {
	as.s.log.Info(ctx, "Process started")
	defer as.s.log.Info(ctx, "Process finished")

	user, err := as.s.accountRepo.GetByID(id)
	if err != nil {
		as.s.log.Errorf(ctx, "error to find user: %s", err.Error())
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (as *accountService) UpdateAccount(ctx context.Context, i dto.UpdateAccount, id string) error {
	as.s.log.Info(ctx, "Process started")
	defer as.s.log.Info(ctx, "Process finished")

	if i.Email != "" {
		_, err := as.s.accountRepo.GetByEmail(i.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			as.s.log.Errorf(ctx, "error to get account by email: %s", err.Error())
			return ErrEmailNotFound
		} else if err == nil {
			as.s.log.Error(ctx, "email already taken")
			return ErrEmailAlreadyTaken
		}
	}

	account := entity.Account{
		Name:  i.Name,
		Email: i.Email,
	}

	err := as.s.accountRepo.Update(&account, id)
	if err != nil {
		as.s.log.Errorf(ctx, "error update user: %s", err.Error())
		return ErrUpdateUser
	}

	return nil
}

func (as *accountService) DeleteAccount(ctx context.Context, id string) error {
	as.s.log.Info(ctx, "Process started")
	defer as.s.log.Info(ctx, "Process finished")

	_, err := as.s.accountRepo.GetByID(id)
	if err != nil {
		as.s.log.Errorf(ctx, "error find user by ID: %s", err.Error())
		return ErrUserNotFound
	}

	err = as.s.accountRepo.Delete(id)
	if err != nil {
		as.s.log.Errorf(ctx, "error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	return nil
}

func (as *accountService) GetAll(ctx context.Context) (*[]entity.Account, error) {
	as.s.log.Info(ctx, "Process started")
	defer as.s.log.Info(ctx, "Process finished")

	accounts, err := as.s.accountRepo.GetAll()
	if err != nil {
		if len(*accounts) == 0 {
			return nil, ErrEmptyResourceError
		}
		as.s.log.Errorf(ctx, "error find users: %s", err.Error())
		return nil, ErrGetManyUsers
	}

	return accounts, nil
}
