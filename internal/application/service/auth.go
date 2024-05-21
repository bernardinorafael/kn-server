package service

import (
	"errors"
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

func (s *authService) Login(email, password string) error {
	_, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) Register(name, email, password string) error {
	// TODO: remove gorm dependency in this checking
	_, err := s.userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	user, err := entity.NewUser(name, email, password)
	if err != nil {
		s.l.Error("error creating user entity", err)
		return err
	}

	err = s.userRepo.Create(*user)
	if err != nil {
		if strings.Contains(err.Error(), "uni_users_email") {
			return ErrEmailAlreadyTaken
		}
		return err
	}

	return nil
}
