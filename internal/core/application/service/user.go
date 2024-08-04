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
	userRepo contract.UserRepository
}

func NewUserService(log logger.Logger, userRepo contract.UserRepository) contract.UserService {
	return &userService{log, userRepo}
}

func (svc userService) RecoverPassword(publicID string, data dto.UpdatePassword) error {
	found, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "id", publicID)
		return ErrUserNotFound
	}

	p, err := password.New(data.NewPassword)
	if err != nil {
		svc.log.Error("error creating password value object", "error", err.Error())
		return err
	}

	err = p.Compare(password.Password(found.Password), data.OldPassword)
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
		Document: found.Document,
		Phone:    found.Phone,
		TeamID:   found.PublicTeamID,
	})
	if err != nil {
		svc.log.Error("failed to initialize new user entity", "error", err.Error())
		return err
	}

	_, err = svc.userRepo.Update(*u)
	if err != nil {
		svc.log.Error("error updating password", "error", err.Error())
		return ErrUpdatingPassword
	}

	return nil
}

func (svc userService) Update(publicID string, data dto.UpdateUser) error {
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
		Document: record.Document,
		Phone:    record.Phone,
		TeamID:   record.PublicTeamID,
	})
	if err != nil {
		svc.log.Error("failed to initialize new user entity", "error", err.Error())
		return err
	}

	if data.Name != "" {
		err = u.ChangeName(data.Name)
		if err != nil {
			svc.log.Error("error while changing name", "error", err.Error())
			return err
		}
	}
	if data.Email != "" {
		err = u.ChangeEmail(data.Email)
		if err != nil {
			svc.log.Error("error while changing email", "error", err.Error())
			return err
		}
	}
	if data.Phone != "" {
		err = u.ChangePhone(data.Phone)
		if err != nil {
			svc.log.Error("error while changing phone", "error", err.Error())
			return err
		}
	}

	_, err = svc.userRepo.Update(*u)
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
