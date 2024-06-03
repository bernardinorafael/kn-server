package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/user"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrInvalidCredential = errors.New("invalid credentials")
)

type authService struct {
	l        *slog.Logger
	userRepo contract.UserRepository
}

func NewAuthService(l *slog.Logger, userRepo contract.UserRepository) contract.AuthService {
	return &authService{l, userRepo}
}

// TODO: implements lockout and rate limiting
func (s *authService) Login(email, password string) (*user.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		s.l.Error(fmt.Sprintf("user with email %s does not exist", email))
		return nil, ErrInvalidCredential
	}

	// err = crypto.Compare(user.Password, password)
	// if err != nil {
	// 	s.l.Error("the password provided is incorrect")
	// 	return nil, ErrInvalidCredential
	// }

	s.l.Info(
		"user attempts to login",
		"name", user.Name, "email", email,
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
		s.l.Error("error creating user entity", err.Error(), err)
		return nil, err
	}

	s.l.Info(
		"registering new user",
		"name", name, "email", email,
	)

	user, err := s.userRepo.Create(*newUser)
	if err != nil {
		if strings.Contains(err.Error(), "uni_users_email") {
			s.l.Error(fmt.Sprintf("email %s is already taken", email))
			return nil, ErrEmailAlreadyTaken
		}
		return nil, err
	}

	return user, nil
}
