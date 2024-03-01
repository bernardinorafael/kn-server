package service

import (
	"errors"
	"time"

	"github.com/bernardinorafael/kn-server/config"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	utillog "github.com/bernardinorafael/kn-server/util/log"
)

var jwtTokenExpiresAt time.Duration

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
	expiredTokenError      = errors.New("the provided access token has expired")
	invalidTokenError      = errors.New("the provided access token is invalid")
	couldNotParseJwtError  = errors.New("failed to parse the provided jwt token")
	encryptTokenError      = errors.New("failed to encrypt the provided token")
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
