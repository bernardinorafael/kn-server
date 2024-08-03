package service

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
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
	if data.Document != "" {
		err = u.ChangeDocument(data.Document)
		if err != nil {
			svc.log.Error("error while changing document", "error", err.Error())
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
