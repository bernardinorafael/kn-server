package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/password"
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

// TODO: implements lockout and rate limiting
func (s *authService) Login(mail, pass string) (*user.User, error) {
	address, err := email.New(mail)
	if err != nil {
		s.log.Error("error creating email value object", err.Error(), err)
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(string(address.ToEmail()))
	if err != nil {
		s.log.Error(fmt.Sprintf("user with email %s does not exist", address.ToEmail()))
		return nil, ErrInvalidCredential
	}

	passw, err := password.New(pass)
	if err != nil {
		s.log.Error("error creating password value object", err.Error(), err)
		return nil, err
	}

	err = passw.Compare(user.Password, pass)
	if err != nil {
		s.log.Error("the password provided is incorrect")
		return nil, ErrInvalidCredential
	}
	return user, nil
}

func (s *authService) Register(name, emailAddr, password, document string) (*user.User, error) {
	nu, err := user.New(name, emailAddr, password, document)
	if err != nil {
		s.log.Error("error creating user entity", err.Error(), err)
		return nil, err
	}

	_, err = s.userRepo.GetByEmail(string(nu.Email))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user, err := s.userRepo.Create(*nu)
	if err != nil {
		if strings.Contains(err.Error(), "uni_users_email") {
			s.log.Error(fmt.Sprintf("email [%s] already exist", string(nu.Email)))
			return nil, ErrEmailAlreadyTaken
		}
		if strings.Contains(err.Error(), "uni_users_document") {
			s.log.Error(fmt.Sprintf("document [%s] already taken", document))
			return nil, ErrDocumentAlreadyTaken
		}
		return nil, err
	}
	return user, nil
}

func (s *authService) RecoverPassword(publicID string, data dto.UpdatePassword) error {
	u, err := s.userRepo.GetByPublicID(publicID)
	if err != nil {
		s.log.Error(fmt.Sprintf("not found user with ID: %d", u.ID))
		return ErrUserNotFound
	}

	pass, err := password.New(data.NewPassword)
	if err != nil {
		s.log.Error("error creating password value object", err.Error(), err)
		return err
	}

	err = pass.Compare(u.Password, data.OldPassword)
	if err != nil {
		s.log.Error("failed to compare password", err.Error(), err)
		return err
	}

	hashed, err := pass.ToEncrypted()
	if err != nil {
		s.log.Error("failed to encrypt password", err.Error(), err)
		return err
	}

	updatedPassword := user.User{PublicID: publicID, Password: hashed}

	_, err = s.userRepo.Update(updatedPassword)
	if err != nil {
		s.log.Error(err.Error())
		return ErrUpdatingPassword
	}
	return nil
}
