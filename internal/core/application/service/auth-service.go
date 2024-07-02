package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/password"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyTaken    = errors.New("email already taken")
	ErrDocumentAlreadyTaken = errors.New("document already taken")
	ErrInvalidCredential    = errors.New("invalid credentials")
	ErrUpdatingPassword     = errors.New("error while updating password")
	ErrUserNotFound         = errors.New("user not found")
)

type authService struct {
	log      *slog.Logger
	userRepo contract.UserRepository
}

func NewAuthService(log *slog.Logger, userRepo contract.UserRepository) contract.AuthService {
	return &authService{log, userRepo}
}

func (svc *authService) Login(userEmail, userPassword string) (*model.User, error) {
	address, err := email.New(userEmail)
	if err != nil {
		svc.log.Error("error creating email value object", err.Error(), err)
		return nil, err
	}

	user, err := svc.userRepo.GetByEmail(string(address.ToEmail()))
	if err != nil {
		svc.log.Error(fmt.Sprintf("user with email %s does not exist", address.ToEmail()))
		return nil, ErrInvalidCredential
	}

	p, err := password.New(userPassword)
	if err != nil {
		svc.log.Error("error creating password value object", err.Error(), err)
		return nil, err
	}

	err = p.Compare(user.Password, userPassword)
	if err != nil {
		svc.log.Error("the password provided is incorrect")
		return nil, ErrInvalidCredential
	}
	return user, nil
}

func (svc *authService) Register(name, emailAddr, password, document string) (*model.User, error) {
	nu, err := user.New(name, emailAddr, password, document)
	if err != nil {
		svc.log.Error("error creating user entity", "error", err.Error())
		return nil, err
	}

	_, err = svc.userRepo.GetByEmail(string(nu.Email))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	u, err := svc.userRepo.Create(*nu)
	if err != nil {
		if strings.Contains(err.Error(), "uni_users_email") {
			svc.log.Error(fmt.Sprintf("email [%s] already exist", string(nu.Email)))
			return nil, ErrEmailAlreadyTaken
		}
		if strings.Contains(err.Error(), "uni_users_document") {
			svc.log.Error(fmt.Sprintf("document [%s] already taken", document))
			return nil, ErrDocumentAlreadyTaken
		}
		return nil, err
	}
	return u, nil
}

func (svc *authService) RecoverPassword(publicID string, data dto.UpdatePassword) error {
	u, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error(fmt.Sprintf("not found user with ID: %d", u.ID))
		return ErrUserNotFound
	}

	p, err := password.New(data.NewPassword)
	if err != nil {
		svc.log.Error("error creating password value object", err.Error(), err)
		return err
	}

	err = p.Compare(u.Password, data.OldPassword)
	if err != nil {
		svc.log.Error("failed to compare password", err.Error(), err)
		return err
	}

	hashed, err := p.ToEncrypted()
	if err != nil {
		svc.log.Error("failed to encrypt password", err.Error(), err)
		return err
	}

	updatedPassword := user.User{PublicID: publicID, Password: hashed}

	_, err = svc.userRepo.Update(updatedPassword)
	if err != nil {
		svc.log.Error(err.Error())
		return ErrUpdatingPassword
	}
	return nil
}
