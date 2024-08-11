package service

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/password"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"github.com/bernardinorafael/kn-server/pkg/logger"
)

type userService struct {
	log           logger.Logger
	userRepo      contract.GormUserRepository
	emailVerifier contract.EmailVerifier
}

func NewUserService(
	log logger.Logger,
	userRepo contract.GormUserRepository,
	emailVerifier contract.EmailVerifier,
) contract.UserService {
	return &userService{
		log:           log,
		userRepo:      userRepo,
		emailVerifier: emailVerifier,
	}
}

func (svc userService) NotifyValidationByEmail(publicID string) error {
	r, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "public_id", publicID)
		return ErrUserNotFound
	}

	u, err := user.New(user.Params{
		PublicID: r.PublicID,
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		Phone:    r.Phone,
		TeamID:   r.PublicTeamID,
	})
	if err != nil {
		svc.log.Error("entity user error", "error", err.Error())
		return err
	}

	err = svc.emailVerifier.NotifyEmail(r.Email)
	if err != nil {
		svc.log.Error("error validating verify code", "error", err.Error())
		return err
	}

	if err = u.ChangeStatus(user.StatusActivationSent); err != nil {
		svc.log.Error("error changing user status", "error", err.Error())
		return err
	}

	userModel := gormodel.User{
		ID:           0,
		PublicID:     u.PublicID(),
		Name:         u.Name(),
		Email:        string(u.Email()),
		Phone:        string(u.Phone()),
		Status:       u.StatusString(),
		Password:     string(u.Password()),
		PublicTeamID: u.TeamID(),
	}

	if _, err = svc.userRepo.Update(userModel); err != nil {
		svc.log.Error("error updating user", "error", err.Error())
		return err
	}

	return nil
}

func (svc userService) ValidateUserByEmail(publicID string, dto dto.ValidateUserByEmail) error {
	found, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "public_id", publicID)
		return ErrUserNotFound
	}

	err = svc.emailVerifier.ConfirmEmail(dto.Code, dto.Email)
	if err != nil {
		svc.log.Error("error validating verify code", "error", err.Error())
		return err
	}

	u, err := user.New(user.Params{
		PublicID: found.PublicID,
		Name:     found.Name,
		Email:    found.Email,
		Password: found.Password,
		Phone:    found.Phone,
		TeamID:   found.PublicTeamID,
	})
	if err != nil {
		svc.log.Error("entity user error", "error", err.Error())
		return err
	}

	if err = u.ChangeStatus(user.StatusEnabled); err != nil {
		svc.log.Error("error changing user status", "error", err.Error())
		return err
	}

	userModel := gormodel.User{
		PublicID:     u.PublicID(),
		Name:         u.Name(),
		Email:        string(u.Email()),
		Phone:        string(u.Phone()),
		Status:       u.StatusString(),
		Password:     string(u.Password()),
		PublicTeamID: u.TeamID(),
	}

	if _, err = svc.userRepo.Update(userModel); err != nil {
		svc.log.Error("error updating user", "error", err.Error())
		return err
	}

	return nil
}

func (svc userService) RecoverPassword(publicID string, dto dto.UpdatePassword) error {
	found, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "id", publicID)
		return ErrUserNotFound
	}

	p, err := password.New(dto.NewPassword)
	if err != nil {
		svc.log.Error("error creating password value object", "error", err.Error())
		return err
	}

	err = p.Compare(password.Password(found.Password), dto.OldPassword)
	if err != nil {
		svc.log.Error("failed to compare password", "error", err.Error())
		return err
	}

	hashed, err := p.Encrypt()
	if err != nil {
		svc.log.Error("failed to encrypt password", "error", err.Error())
		return err
	}

	u, err := user.New(user.Params{
		PublicID: found.PublicID,
		Name:     found.Name,
		Email:    found.Email,
		Password: string(hashed),
		Phone:    found.Phone,
		TeamID:   found.PublicTeamID,
	})
	if err != nil {
		svc.log.Error("failed to initialize new user entity", "error", err.Error())
		return err
	}

	userModel := gormodel.User{
		PublicID:     u.PublicID(),
		Name:         u.Name(),
		Email:        string(u.Email()),
		Phone:        string(u.Phone()),
		Status:       u.StatusString(),
		Password:     string(u.Password()),
		PublicTeamID: u.TeamID(),
	}

	if _, err = svc.userRepo.Update(userModel); err != nil {
		svc.log.Error("error updating user", "error", err.Error())
		return err
	}

	return nil
}

func (svc userService) Update(publicID string, dto dto.UpdateUser) error {
	found, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "public_id", publicID)
		return errors.New("user not found")
	}

	u, err := user.New(user.Params{
		PublicID: found.PublicID,
		Name:     found.Name,
		Email:    found.Email,
		Password: found.Password,
		Phone:    found.Phone,
		TeamID:   found.PublicTeamID,
	})
	if err != nil {
		svc.log.Error("failed to initialize new user entity", "error", err.Error())
		return err
	}

	if dto.Name != "" {
		if err = u.ChangeName(dto.Name); err != nil {
			svc.log.Error("error while changing name", "error", err.Error())
			return err
		}
	}
	if dto.Email != "" {
		if err = u.ChangeEmail(dto.Email); err != nil {
			svc.log.Error("error while changing email", "error", err.Error())
			return err
		}
	}
	if dto.Phone != "" {
		if err = u.ChangePhone(dto.Phone); err != nil {
			svc.log.Error("error while changing phone", "error", err.Error())
			return err
		}
	}

	userModel := gormodel.User{
		PublicID:     u.PublicID(),
		Name:         u.Name(),
		Email:        string(u.Email()),
		Phone:        string(u.Phone()),
		Status:       u.StatusString(),
		Password:     string(u.Password()),
		PublicTeamID: u.TeamID(),
	}

	if _, err = svc.userRepo.Update(userModel); err != nil {
		svc.log.Error("error updating user", "error", err.Error())
		return err
	}

	return nil
}

func (svc userService) GetUser(publicID string) (gormodel.User, error) {
	user, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "public_id", publicID)
		return user, errors.New("user not found")
	}

	return user, nil
}
