package service

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	utillog "github.com/bernardinorafael/kn-server/util/log"
)

var (
	userAlreadyTakenError  = errors.New("the provided user is already in use")
	hashPasswordError      = errors.New("an error occurred in trying to hash password")
	createUserError        = errors.New("an error occurred trying to create user")
	userNotFoundError      = errors.New("no users were found with the provided ID")
	emailNotFoundError     = errors.New("no users matched the provided e-mail")
	updateUserError        = errors.New("an error occurred, cannot update this resource")
	deleteUserError        = errors.New("an error occurred, cannot delete this resource")
	getManyUsersError      = errors.New("an error occurred, unable to retrieve the resource")
	invalidCredentialError = errors.New("the provided input does not match the server")
	equalPasswordsError    = errors.New("both passwords are the same")
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
