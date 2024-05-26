package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyTaken = errors.New("email already taken")
)

type authService struct {
	l        *slog.Logger
	userRepo contract.UserRepository
}

func NewAuthService(l *slog.Logger, userRepo contract.UserRepository) contract.AuthService {
	return &authService{l, userRepo}
}

func (s *authService) Login(email, password string) (*entity.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Register(name, email, password string) (*entity.User, error) {
	_, err := s.userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	newUser, err := entity.NewUser(name, email, password)
	if err != nil {
		s.l.Error("error creating user entity", err)
		return nil, err
	}

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
