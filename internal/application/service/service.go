package service

import (
	"errors"

	"github.com/bernardinorafael/kn-server/config"
	utillog "github.com/bernardinorafael/kn-server/helper/log"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
)

var (
	ErrUserAlreadyTaken     = errors.New("user already taken")
	ErrDocumentALreadyTaken = errors.New("the provided document is already registered in the system")
	ErrEmailAlreadyTaken    = errors.New("email already taken")
	ErrHashPassword         = errors.New("an error occurred in trying to hash password")
	ErrCreateUser           = errors.New("an error occurred trying to create user")
	ErrUserNotFound         = errors.New("no users were found with the provided ID")
	ErrEmailNotFound        = errors.New("no users were found with the provided email address")
	ErrUpdateUser           = errors.New("an error occurred, cannot update this resource")
	ErrDeleteUser           = errors.New("an error occurred, cannot delete this resource")
	ErrGetManyUsers         = errors.New("an error occurred, unable to retrieve the resource")
	ErrInvalidCredentials   = errors.New("authentication failed, the provided email and/or password is incorrect")
	ErrEqualPasswords       = errors.New("both passwords are the same")
	ErrExpiredToken         = errors.New("the provided access token has expired")
	ErrCouldNotParseJWT     = errors.New("failed to parse the provided jwt token")
	ErrEncryptToken         = errors.New("failed to encrypt the provided token")
	ErrEmptyResourceError   = errors.New("there is not resource yet, the list is empty")
)

type service struct {
	accountRepo contract.AccountRepository
	log         utillog.Logger
	cfg         *config.EnvFile
}

type Services struct {
	AccountService contract.AccountService
	JWTService     contract.JWTService
	AuthService    contract.AuthService
}

type svcOptions func(*service)

func New(svcOptions ...svcOptions) *Services {
	svc := &service{}
	for _, opt := range svcOptions {
		opt(svc)
	}

	return &Services{
		AccountService: newAccountService(svc),
		JWTService:     newJWTService(svc),
		AuthService:    newAuthService(svc),
	}
}

func GetAccountRepository(accountRepo contract.AccountRepository) svcOptions {
	return func(service *service) {
		service.accountRepo = accountRepo
	}
}
func GetConfig(cfg *config.EnvFile) svcOptions {
	return func(service *service) {
		service.cfg = cfg
	}
}
func GetLogger(log utillog.Logger) svcOptions {
	return func(service *service) {
		service.log = log
	}
}
