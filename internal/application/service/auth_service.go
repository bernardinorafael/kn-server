package service

import (
	"errors"
	"strings"
	"time"

	"github.com/bernardinorafael/kn-server/helper/crypto"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authService struct {
	s *service
}

func newAuthService(service *service) contract.AuthService {
	return &authService{
		s: service,
	}
}

func (us *authService) Register(input dto.Register) (*entity.User, error) {
	us.s.log.Info("Process started")
	defer us.s.log.Info("Process finished")

	_, err := us.s.userRepo.GetByEmail(input.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		us.s.log.Error("error to get account by email", err)
		return nil, ErrEmailNotFound
	} else if err == nil {
		us.s.log.Error("error to get account by email", ErrEmailAlreadyTaken)
		return nil, ErrEmailAlreadyTaken
	}

	encrypted, err := crypto.EncryptPassword(input.Password)
	if err != nil {
		us.s.log.Error("failed to encrypt password", err)
		return nil, ErrEncryptToken
	}

	user := entity.User{
		ID:        uuid.New().String(),
		Name:      input.Name,
		Surname:   input.Surname,
		Email:     input.Email,
		Password:  encrypted,
		Document:  input.Document,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = us.s.userRepo.Save(user)
	if err != nil {
		if strings.Contains(err.Error(), "uni_users_document") {
			us.s.log.Error("document already exist", ErrDocumentAlreadyTaken)
			return nil, ErrDocumentAlreadyTaken
		}
		us.s.log.Error("error creating account", err)
		return nil, ErrCreateUser
	}

	us.s.log.Info("creating account to",
		"name", user.Name,
		"surname", user.Surname,
		"email", user.Email,
	)

	return &user, nil
}

func (us *authService) Login(input dto.Login) (*entity.User, error) {
	us.s.log.Info("Process started")
	defer us.s.log.Info("Process finished")

	user, err := us.s.userRepo.GetByEmail(input.Email)
	if err != nil {
		us.s.log.Error("cannot find user by email", err)
		return nil, ErrInvalidCredentials
	}

	encrypted, err := us.s.userRepo.GetPassword(user.ID)
	if err != nil {
		us.s.log.Error("error to get user password", err)
		return nil, ErrInvalidCredentials
	}

	err = crypto.CheckPassword(encrypted, input.Password)
	if err != nil {
		us.s.log.Error("password does not match", err)
		return nil, ErrInvalidCredentials
	}

	us.s.log.Info("login process to user",
		"email", user.Email,
	)

	return user, nil
}
