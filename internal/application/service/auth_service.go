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

func (us *authService) Register(ctx context.Context, i dto.Register) (*entity.Account, error) {
	us.s.log.Info(ctx, "Process started")
	defer us.s.log.Info(ctx, "Process finished")

	_, err := us.s.accountRepo.GetByEmail(i.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		us.s.log.Errorf(ctx, "error to get account by email: %s", err.Error())
		return nil, ErrEmailNotFound
	} else if err == nil {
		return nil, ErrEmailAlreadyTaken
	}

	encrypted, err := crypto.EncryptPassword(i.Password)
	if err != nil {
		us.s.log.Error(ctx, "failed to encrypt password")
		return nil, ErrEncryptToken
	}

	account := entity.Account{
		ID:        uuid.New().String(),
		Name:      i.Name,
		Email:     i.Email,
		Password:  encrypted,
		Document:  i.Document,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = us.s.accountRepo.Save(account)
	if err != nil {
		if strings.Contains(err.Error(), "uni_accounts_document") {
			us.s.log.Error(ctx, "document already exist")
			return nil, ErrDocumentAlreadyTaken
		}
		us.s.log.Error(ctx, "error creating account")
		return nil, ErrCreateUser
	}

	ctx = context.WithValue(ctx, restutil.AuthKey, account.ID)

	return &account, nil
}

func (us *authService) Login(ctx context.Context, i dto.Login) (*entity.Account, error) {
	us.s.log.Info(ctx, "Process started")
	defer us.s.log.Info(ctx, "Process finished")

	acc, err := us.s.accountRepo.GetByEmail(i.Email)
	if err != nil {
		us.s.log.Errorf(ctx, "cannot find user by email: %s", err.Error())
		return nil, ErrInvalidCredentials
	}

	encrypted, err := us.s.accountRepo.GetPassword(acc.ID)
	if err != nil {
		us.s.log.Errorf(ctx, "error to get user password: %s", err.Error())
		return nil, ErrInvalidCredentials
	}

	err = crypto.CheckPassword(encrypted, i.Password)
	if err != nil {
		us.s.log.Errorf(ctx, "password does not match: %s", err.Error())
		return nil, ErrInvalidCredentials
	}

	ctx = context.WithValue(ctx, restutil.AuthKey, acc.ID)

	return acc, nil
}
