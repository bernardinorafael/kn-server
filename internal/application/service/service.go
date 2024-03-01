package service

import (
	"errors"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	utillog "github.com/bernardinorafael/kn-server/util/log"
)

var (
	errUserAlreadyTaken  = errors.New("the provided user is already in use")
	errHashPassword      = errors.New("an error occurred in trying to hash password")
	errCreateUser        = errors.New("an error occurred trying to create user")
	errUserNotFound      = errors.New("no users were found with the provided ID")
	errEmailNotFound     = errors.New("no users matched the provided e-mail")
	errUpdateUser        = errors.New("an error occurred, cannot update this resource")
	errDeleteUser        = errors.New("an error occurred, cannot delete this resource")
	errGetManyUsers      = errors.New("an error occurred, unable to retrieve the resource")
	errInvAlidCredential = errors.New("the provided input does not match the server")
	errEqualPasswords    = errors.New("both passwords are the same")
	errExpiredToken      = errors.New("the provided access token has expired")
	errInvalidToken      = errors.New("the provided access token is invalid")
	errCouldNotParseJwt  = errors.New("failed to parse the provided jwt token")
	errEncryptToken      = errors.New("failed to encrypt the provided token")
)

type service struct {
	ar  contract.AccountRepository
	l   utillog.Logger
	cfg *config.EnvFile
}

type Services struct {
	AccountService contract.AccountService
	AuthService    contract.AuthService
}

type svcOptions func(*service)

func New(svcOptions ...svcOptions) (*Services, error) {
	svc := &service{}
	for _, opt := range svcOptions {
		opt(svc)
	}
	return &Services{
		AccountService: newAccountService(svc),
		AuthService:    newAuthService(svc),
	}, nil
}

func GetAccountRepository(ar contract.AccountRepository) svcOptions {
	return func(s *service) {
		s.ar = ar
	}
}
func GetConfig(cfg *config.EnvFile) svcOptions {
	return func(s *service) {
		s.cfg = cfg
	}
}
func GetLogger(l utillog.Logger) svcOptions {
	return func(s *service) {
		s.l = l
	}
}
