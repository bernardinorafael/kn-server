package service

import (
	"context"
	"errors"
	"log/slog"
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
	us.s.log.Info("Process started")
	defer us.s.log.Info("Process finished")

	_, err := us.s.accountRepo.GetByEmail(i.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		us.s.log.Error("error to get account by email", err.Error())
		return nil, ErrEmailNotFound
	} else if err == nil {
		return nil, ErrEmailAlreadyTaken
	}

	encrypted, err := crypto.EncryptPassword(i.Password)
	if err != nil {
		us.s.log.Error("failed to encrypt password")
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
			us.s.log.Error("document already exist")
			return nil, ErrDocumentAlreadyTaken
		}
		us.s.log.Error("error creating account")
		return nil, ErrCreateUser
	}

	ctx = context.WithValue(ctx, restutil.AuthKey, account.ID)

	return &account, nil
}

func (us *authService) Login(ctx context.Context, i dto.Login) (*entity.Account, error) {
	us.s.log.Info("Process started")
	defer us.s.log.Info("Process finished")

	acc, err := us.s.accountRepo.GetByEmail(i.Email)
	if err != nil {
		us.s.log.Error("cannot find user by email", err.Error())
		return nil, ErrInvalidCredentials
	}

	us.s.log.With(slog.Group("user trying to login",
		"name", acc.Name, "email", acc.Email,
	))

	encrypted, err := us.s.accountRepo.GetPassword(acc.ID)
	if err != nil {
		us.s.log.Error("error to get user password", err.Error())
		return nil, ErrInvalidCredentials
	}

	err = crypto.CheckPassword(encrypted, i.Password)
	if err != nil {
		us.s.log.Error("password does not match", err.Error())
		return nil, ErrInvalidCredentials
	}

	ctx = context.WithValue(ctx, restutil.AuthKey, acc.ID)

	return acc, nil
}
