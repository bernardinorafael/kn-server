package service

import "github.com/bernardinorafael/gozinho/internal/application/contract"

type service struct {
	ar contract.AccountRepository
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
