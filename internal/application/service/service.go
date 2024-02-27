package service

import (
	"github.com/bernardinorafael/gozinho/internal/application/contract"
	utillog "github.com/bernardinorafael/gozinho/util/log"
)

type service struct {
	ar contract.AccountRepository
	l  utillog.Logger
}

type Services struct {
	AccountService contract.AccountService
}

type Options func(*service)

func New(opts ...Options) (*Services, error) {
	svc := &service{}
	for _, option := range opts {
		option(svc)
	}

	return &Services{
		AccountService: newAccountService(svc),
	}, nil
}

func GetAccountRepository(ar contract.AccountRepository) Options {
	return func(service *service) {
		service.ar = ar
	}
}

func GetLogger(l utillog.Logger) Options {
	return func(service *service) {
		service.l = l
	}
}
