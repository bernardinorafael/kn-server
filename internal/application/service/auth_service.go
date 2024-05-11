package service

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
)

type authService struct {
	svc *service
}

func newAuthService(svc *service) contract.AuthService {
	return &authService{svc}
}

func (s authService) Login(email, password string) error {
	s.svc.log.Info("Start process")
	defer s.svc.log.Info("End process")

	return nil
}

func (s authService) Register(name, email, password string) error {
	s.svc.log.Info("Start process")
	defer s.svc.log.Info("End process")

	return nil
}
