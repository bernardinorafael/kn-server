package service

import (
	"log/slog"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type authService struct {
	log      *slog.Logger
	userRepo contract.UserRepository
}

func NewAuthService(userRepo contract.UserRepository, log *slog.Logger) contract.AuthService {
	return &authService{
		log:      log,
		userRepo: userRepo,
	}
}

func (s *authService) Login(email, password string) error {
	s.log.Info("Start process")
	defer s.log.Info("End process")

	_, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) Register(name, email, password string) error {
	s.log.Info("Start process")
	defer s.log.Info("End process")

	_, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	user, err := entity.NewUser(name, email, password)
	if err != nil {
		return err
	}

	if err = s.userRepo.Create(*user); err != nil {
		return err
	}

	return nil
}
