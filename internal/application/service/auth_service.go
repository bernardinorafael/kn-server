package service

import (
	"context"
	"errors"
	"time"

	"github.com/bernardinorafael/kn-server/helper/crypto"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authService struct {
	s *service
}

func newAuthService(service *service) contract.AuthService {
	return &authService{service}
}

func (us *authService) Register(ctx context.Context, i dto.Register) (id string, err error) {
	us.s.log.Info(ctx, "Process started")
	defer us.s.log.Info(ctx, "Process finished")

	_, err = us.s.accountRepo.GetByEmail(i.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		us.s.log.Errorf(ctx, "error to get account by email: %s", err.Error())
		return "", ErrEmailNotFound
	} else if err == nil {
		return "", ErrEmailAlreadyTaken
	}

	encrypted, err := crypto.EncryptPassword(i.Password)
	if err != nil {
		us.s.log.Error(ctx, "failed to encrypt password")
		return "", ErrEncryptToken
	}

	user := entity.Account{
		ID:        uuid.New().String(),
		Name:      i.Name,
		Email:     i.Email,
		Password:  encrypted,
		Document:  i.Document,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = us.s.accountRepo.Save(user)
	if err != nil {
		us.s.log.Errorf(ctx, "error creating user: %s", err.Error())
		return "", ErrCreateUser
	}

	ctx = context.WithValue(ctx, restutil.AuthKey, user.ID)

	return user.ID, nil
}

func (us *authService) Login(ctx context.Context, i dto.Login) (id string, err error) {
	us.s.log.Info(ctx, "Process started")
	defer us.s.log.Info(ctx, "Process finished")

	account, err := us.s.accountRepo.GetByEmail(i.Email)
	if err != nil {
		us.s.log.Errorf(ctx, "cannot find user by email: %s", err.Error())
		return "", ErrInvalidCredentials
	}

	encrypted, err := us.s.accountRepo.GetPassword(account.ID)
	if err != nil {
		us.s.log.Errorf(ctx, "error to get user password: %s", err.Error())
		return "", ErrInvalidCredentials
	}

	err = crypto.CheckPassword(encrypted, i.Password)
	if err != nil {
		us.s.log.Errorf(ctx, "password does not match: %s", err.Error())
		return "", ErrInvalidCredentials
	}

	ctx = context.WithValue(ctx, restutil.AuthKey, account.ID)

	return account.ID, nil
}
