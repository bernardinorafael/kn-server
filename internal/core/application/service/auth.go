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
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrPhoneAlreadyTaken = errors.New("phone already taken")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrUpdatingPassword  = errors.New("error while updating password")
	ErrUserNotFound      = errors.New("user not found")
)

type authService struct {
	log         logger.Logger
	smsVerifier contract.SMSVerifier
	userRepo    contract.GormUserRepository
}

func NewAuthService(
	log logger.Logger,
	smsVerifier contract.SMSVerifier,
	userRepo contract.GormUserRepository,
) contract.AuthService {
	return &authService{
		log:         log,
		smsVerifier: smsVerifier,
		userRepo:    userRepo,
	}
}

func (svc authService) NotifyLoginOTP(dto dto.NotifySMS) error {
	if _, err := phone.New(dto.Phone); err != nil {
		svc.log.Error("failed to validate phone", "error", err.Error())
		return err
	}

	user, err := svc.userRepo.GetByPhone(dto.Phone)
	if err != nil {
		svc.log.Error("not found user by phone", "phone", dto.Phone)
		return ErrUserNotFound
	}

	if err = svc.smsVerifier.NotifySMS(user.Phone); err != nil {
		svc.log.Error("code verification failed", "error", err.Error())
		return err
	}

	return nil
}

func (svc authService) LoginOTP(dto dto.LoginOTP) (gormodel.User, error) {
	var user gormodel.User

	if _, err := phone.New(dto.Phone); err != nil {
		svc.log.Error("failed to validate phone", "error", err.Error())
		return user, err
	}

	user, err := svc.userRepo.GetByPhone(dto.Phone)
	if err != nil {
		svc.log.Error("not found user by phone", "phone", dto.Phone)
		return user, ErrUserNotFound
	}

	if err = svc.smsVerifier.ConfirmSMS(dto.Code, dto.Phone); err != nil {
		svc.log.Error("code verification failed", "error", err.Error())
		return user, err
	}

	return user, nil
}

func (svc authService) Login(dto dto.Login) (gormodel.User, error) {
	var user gormodel.User

	if _, err := email.New(dto.Email); err != nil {
		svc.log.Error("error creating email value object", "error", err.Error())
		return user, err
	}

	user, err := svc.userRepo.GetByEmail(dto.Email)
	if err != nil {
		svc.log.Error("failed to find user by email", "email", dto.Email)
		return user, ErrInvalidCredential
	}

	p, err := password.New(dto.Password)
	if err != nil {
		svc.log.Error("error creating password value object", "error", err.Error())
		return user, err
	}

	err = p.Compare(password.Password(user.Password), dto.Password)
	if err != nil {
		svc.log.Error("the password provided is incorrect", "password", dto.Password)
		return user, ErrInvalidCredential
	}

	return user, nil
}

func (svc authService) Register(dto dto.Register) error {
	u, err := user.New(user.Params{
		PublicID: uuid.NewString(),
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Phone:    dto.Phone,
		TeamID:   nil,
	})
	if err != nil {
		svc.log.Error("failed to initialize new user entity", "error", err.Error())
		return err
	}

	if err = u.EncryptPassword(); err != nil {
		svc.log.Error("failed to encrypt password", "error", err.Error())
		return err
	}

	userModel := gormodel.User{
		PublicID: u.PublicID(),
		Name:     u.Name(),
		Email:    string(u.Email()),
		Phone:    string(u.Phone()),
		Status:   u.StatusString(),
		Password: string(u.Password()),
	}

	if err = svc.userRepo.Create(userModel); err != nil {
		if strings.Contains(err.Error(), "uni_users_email") {
			svc.log.Error("email already taken", "email", dto.Email)
			return ErrEmailAlreadyTaken
		}
		if strings.Contains(err.Error(), "uni_users_phone") {
			svc.log.Error("phone already taken", "phone", dto.Phone)
			return ErrPhoneAlreadyTaken
		}
		return err
	}

	return nil
}
