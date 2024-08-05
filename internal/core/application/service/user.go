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
	log      logger.Logger
	userRepo contract.GormUserRepository
}

func NewUserService(log logger.Logger, userRepo contract.GormUserRepository) contract.UserService {
	return &userService{log, userRepo}
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

	hashed, err := p.ToEncrypted()
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

	_, err = svc.userRepo.Update(userModel)
	if err != nil {
		svc.log.Error("error updating password", "error", err.Error())
		return ErrUpdatingPassword
	}

	return nil
}

func (svc userService) Update(publicID string, dto dto.UpdateUser) error {
	record, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "public_id", publicID)
		return errors.New("user not found")
	}

	u, err := user.New(user.Params{
		PublicID: record.PublicID,
		Name:     record.Name,
		Email:    record.Email,
		Password: record.Password,
		Phone:    record.Phone,
		TeamID:   record.PublicTeamID,
	})
	if err != nil {
		svc.log.Error("failed to initialize new user entity", "error", err.Error())
		return err
	}

	if dto.Name != "" {
		err = u.ChangeName(dto.Name)
		if err != nil {
			svc.log.Error("error while changing name", "error", err.Error())
			return err
		}
	}
	if dto.Email != "" {
		err = u.ChangeEmail(dto.Email)
		if err != nil {
			svc.log.Error("error while changing email", "error", err.Error())
			return err
		}
	}
	if dto.Phone != "" {
		err = u.ChangePhone(dto.Phone)
		if err != nil {
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

	_, err = svc.userRepo.Update(userModel)
	if err != nil {
		svc.log.Error("cannot update user", "error", err.Error())
		return err
	}

	return nil
}

func (svc userService) GetUser(publicID string) (gormodel.User, error) {
	user, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "public_id", publicID)
		return gormodel.User{}, errors.New("user not found")
	}

	return user, nil
}
