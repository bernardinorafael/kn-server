package service

import (
	"errors"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/password"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/phone"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyTaken    = errors.New("email already taken")
	ErrDocumentAlreadyTaken = errors.New("document already taken")
	ErrPhoneAlreadyTaken    = errors.New("phone already taken")
	ErrInvalidCredential    = errors.New("invalid credentials")
	ErrUpdatingPassword     = errors.New("error while updating password")
	ErrUserNotFound         = errors.New("user not found")
)

type authService struct {
	log      logger.Logger
	userRepo contract.UserRepository
}

func NewAuthService(log logger.Logger, userRepo contract.UserRepository) contract.AuthService {
	return &authService{log, userRepo}
}

func (svc authService) LoginOTP(dto dto.LoginOTP) (gormodel.User, error) {
	var userModel gormodel.User

	_, err := phone.New(dto.Phone)
	if err != nil {
		svc.log.Error("error creating phone value object", "error", err.Error())
		return userModel, err
	}

	user, err := svc.userRepo.GetByPhone(dto.Phone)
	if err != nil {
		svc.log.Error("failed to find user by phone", "phone", dto.Phone)
		return userModel, ErrUserNotFound
	}

	return user, nil
}

func (svc authService) Login(dto dto.Login) (gormodel.User, error) {
	var userModel gormodel.User

	_, err := email.New(dto.Email)
	if err != nil {
		svc.log.Error("error creating email value object", "error", err.Error())
		return userModel, err
	}

	user, err := svc.userRepo.GetByEmail(dto.Email)
	if err != nil {
		svc.log.Error("failed to find user by email", "email", dto.Email)
		return userModel, ErrInvalidCredential
	}

	p, err := password.New(dto.Password)
	if err != nil {
		svc.log.Error("error creating password value object", "error", err.Error())
		return userModel, err
	}

	err = p.Compare(password.Password(user.Password), dto.Password)
	if err != nil {
		svc.log.Error("the password provided is incorrect", "password", dto.Password)
		return userModel, ErrInvalidCredential
	}

	return user, nil
}

func (svc authService) Register(dto dto.Register) (gormodel.User, error) {
	var userModel gormodel.User

	u, err := user.New(user.Params{
		PublicID: uuid.NewString(),
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Document: dto.Document,
		Phone:    dto.Phone,
		TeamID:   nil,
	})
	if err != nil {
		svc.log.Error("failed to initialize new user entity", "error", err.Error())
		return userModel, err
	}

	if err = u.EncryptPassword(); err != nil {
		svc.log.Error("failed to encrypt password", "error", err.Error())
		return userModel, err
	}

	newUser, err := svc.userRepo.Create(*u)
	if err != nil {
		if strings.Contains(err.Error(), "uni_users_email") {
			svc.log.Error("email already taken", "email", dto.Email)
			return userModel, ErrEmailAlreadyTaken
		}
		if strings.Contains(err.Error(), "uni_users_phone") {
			svc.log.Error("phone already taken", "phone", dto.Phone)
			return userModel, ErrPhoneAlreadyTaken
		}
		if strings.Contains(err.Error(), "uni_users_document") {
			svc.log.Error("document already taken", "document", dto.Document)
			return userModel, ErrDocumentAlreadyTaken
		}
		return userModel, err
	}

	return newUser, nil
}
