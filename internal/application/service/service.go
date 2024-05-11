package service

import (
	"log/slog"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
)

type service struct {
	log slog.Logger
}

type Services struct {
	AuthService contract.AuthService
}

type Options = func(*service)

func New(svcOptions ...Options) *Services {
	svc := &service{}
	for _, opt := range svcOptions {
		opt(svc)
	}

	return &Services{
		AuthService: newAuthService(svc),
	}
}

func GetLogger(log *slog.Logger) Options {
	return func(service *service) {
		service.log = *log
	}
}
