package service

import (
	"errors"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/password"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyTaken    = errors.New("email already taken")
	ErrDocumentAlreadyTaken = errors.New("document already taken")
	ErrPhoneAlreadyTaken    = errors.New("phone already taken")
	ErrInvalidCredential    = errors.New("invalid credentials")
	ErrUpdatingPassword     = errors.New("error while updating password")
	ErrUserNotFound         = errors.New("user not found")
)

type authService struct {
	log      logger.Logger
	userRepo contract.UserRepository
}

func NewAuthService(log logger.Logger, userRepo contract.UserRepository) contract.AuthService {
	return &authService{log, userRepo}
}

func (svc *authService) Login(data dto.Login) (*gormodel.User, error) {
	address, err := email.New(data.Email)
	if err != nil {
		svc.log.Error("error creating email value object", "error", err.Error())
		return nil, err
	}

	user, err := svc.userRepo.GetByEmail(string(address.ToEmail()))
	if err != nil {
		svc.log.Error("cannot find user by the given email", "email", data.Email)
		return nil, ErrInvalidCredential
	}

	p, err := password.New(data.Password)
	if err != nil {
		svc.log.Error("error creating password value object", "error", err.Error())
		return nil, err
	}

	err = p.Compare(password.Password(user.Password), data.Password)
	if err != nil {
		svc.log.Error("the password provided is incorrect", "password", data.Password)
		return nil, ErrInvalidCredential
	}
	return user, nil
}

func (svc *authService) Register(data dto.Register) (*gormodel.User, error) {
	u, err := user.New(user.Params{
		PublicID: uuid.NewString(),
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Document: data.Document,
		Phone:    data.Phone,
		TeamID:   nil,
	})
	if err != nil {
		svc.log.Error("failed to initialize new user entity", "error", err.Error())
		return nil, err
	}

	if err = u.EncryptPassword(); err != nil {
		svc.log.Error("failed to encrypt password", "error", err.Error())
		return nil, err
	}

	newUser, err := svc.userRepo.Create(*u)
	if err != nil {
		if strings.Contains(err.Error(), "uni_users_email") {
			svc.log.Error("email already taken", "email", data.Email)
			return nil, ErrEmailAlreadyTaken
		}
		if strings.Contains(err.Error(), "uni_users_phone") {
			svc.log.Error("phone already taken", "phone", data.Phone)
			return nil, ErrPhoneAlreadyTaken
		}
		if strings.Contains(err.Error(), "uni_users_document") {
			svc.log.Error("document already taken", "document", data.Document)
			return nil, ErrDocumentAlreadyTaken
		}
		return nil, err
	}
	return newUser, nil
}

func (svc *authService) RecoverPassword(publicID string, data dto.UpdatePassword) error {
	record, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "id", publicID)
		return ErrUserNotFound
	}

	p, err := password.New(data.NewPassword)
	if err != nil {
		svc.log.Error("error creating password value object", "error", err.Error())
		return err
	}

	err = p.Compare(password.Password(record.Password), data.OldPassword)
	if err != nil {
		svc.log.Error("failed to compare password", "error", err.Error())
		return err
	}

	hashed, err := p.ToEncrypted()
	if err != nil {
		svc.log.Error("failed to encrypt password", "error", err.Error())
		return err
	}

	u, err := user.New(user.Params{
		PublicID: record.PublicID,
		Name:     record.Name,
		Email:    record.Email,
		Password: string(hashed),
		Document: record.Document,
		Phone:    record.Phone,
		TeamID:   record.PublicTeamID,
	})
	if err != nil {
		svc.log.Error("failed to initialize new user entity", "error", err.Error())
		return err
	}

	_, err = svc.userRepo.Update(*u)
	if err != nil {
		svc.log.Error("error updating password", "error", err.Error())
		return ErrUpdatingPassword
	}
	return nil
}
