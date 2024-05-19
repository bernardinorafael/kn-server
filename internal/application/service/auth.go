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

type authService struct {
	l        *slog.Logger
	userRepo contract.UserRepository
}

func NewAuthService(l *slog.Logger, userRepo contract.UserRepository) contract.AuthService {
	return &authService{l, userRepo}
}

func (s *authService) Login(email, password string) error {
	s.l.Info("Start service process")
	defer s.l.Info("End service process")

	_, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) Register(name, email, password string) error {
	s.l.Info("Start service process")
	defer s.l.Info("End service process")

	// TODO: remove gorm dependency in this checking
	_, err := s.userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	user, err := entity.NewUser(name, email, password)
	if err != nil {
		s.l.Error("error creating user", err)
		return err
	}

	// TODO: verify best way to check sql error
	if err = s.userRepo.Create(*user); err != nil {
		if strings.Contains(err.Error(), "uni_users_email") {
			return fmt.Errorf("email %v is already taken", email)
		}
		return err
	}

	return nil
}
