package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/password"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrUserNotFound      = errors.New("user not found")
	ErrUpdatingPassword  = errors.New("error while updating password")
)

type authService struct {
	log      *slog.Logger
	userRepo contract.UserRepository
}

func NewAuthService(log *slog.Logger, userRepo contract.UserRepository) contract.AuthService {
	return &authService{log, userRepo}
}

// TODO: implements lockout and rate limiting
func (s *authService) Login(email, pass string) (*user.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		s.log.Error(fmt.Sprintf("user with email %s does not exist", email))
		return nil, ErrInvalidCredential
	}

	passw, err := password.New(pass)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	err = passw.Compare(user.Password, pass)
	if err != nil {
		s.log.Error("the password provided is incorrect")
		return nil, ErrInvalidCredential
	}

	s.log.Info(
		"login user",
		"name", user.Name,
		"email", email,
	)

	return user, nil
}

func (s *authService) Register(name, email, password string) (*user.User, error) {
	_, err := s.userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	newUser, err := user.New(name, email, password)
	if err != nil {
		s.log.Error("error creating user entity", err.Error(), err)
		return nil, err
	}

	s.log.Info(
		"creating user",
		"name", name,
		"email", email,
	)

	user, err := s.userRepo.Create(*newUser)
	if err != nil {
		if strings.Contains(err.Error(), "uni_users_email") {
			s.log.Error(fmt.Sprintf("email %s is already taken", email))
			return nil, ErrEmailAlreadyTaken
		}
		return nil, err
	}

	return user, nil
}

func (s *authService) RecoverPassword(id uint, data dto.UpdatePassword) error {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		s.log.Error(fmt.Sprintf("not found user with ID: %d", u.ID))
		return ErrUserNotFound
	}

	pass, err := password.New(data.NewPassword)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	err = pass.Compare(u.Password, data.OldPassword)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	hashed, err := pass.ToEncrypted()
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	updatedPassword := user.User{Model: gorm.Model{ID: id}, Password: hashed}

	_, err = s.userRepo.Update(updatedPassword)
	if err != nil {
		return ErrUpdatingPassword
	}

	return nil
}
