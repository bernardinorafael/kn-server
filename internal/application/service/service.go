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
	AuthService    contract.AccountService
}

type SvcOptions func(*service)

func New(opts ...SvcOptions) (*Services, error) {
	svc := &service{}
	for _, option := range opts {
		option(svc)
	}

	return &Services{AccountService: newAccountService(svc)}, nil
}

func GetAccountRepository(ar contract.AccountRepository) SvcOptions {
	return func(service *service) {
		service.ar = ar
	}
}

func GetLogger(l utillog.Logger) SvcOptions {
	return func(service *service) {
		service.l = l
	}
}
