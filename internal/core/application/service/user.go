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

func (svc *userService) Update(publicID string, data dto.UpdateUser) error {
	foundUser, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "public_id", publicID)
		return errors.New("user not found")
	}

	userEntity, err := user.New(
		foundUser.Name,
		string(foundUser.Email),
		string(foundUser.Password),
		string(foundUser.Document),
		string(foundUser.Phone),
		nil,
	)
	if err != nil {
		svc.log.Error("error init user entity", "error", err.Error())
		return err
	}

	if data.Name != "" {
		err = userEntity.ChangeName(data.Name)
		if err != nil {
			svc.log.Error("erro while changing name", "error", err.Error())
			return err
		}
	}

	if data.Document != "" {
		err = userEntity.ChangeDocument(data.Document)
		if err != nil {
			svc.log.Error("error while changing document", "error", err.Error())
			return err
		}
	}

	if data.Email != "" {
		err = userEntity.ChangeEmail(data.Email)
		if err != nil {
			svc.log.Error("error while changing email", "error", err.Error())
			return err
		}
	}

	if data.Phone != "" {
		err = userEntity.ChangePhone(data.Phone)
		if err != nil {
			svc.log.Error("error while changing phone", "error", err.Error())
			return err
		}
	}

	updated := user.User{
		PublicID: publicID,
		Name:     userEntity.Name,
		Email:    userEntity.Email,
		Document: userEntity.Document,
		Phone:    userEntity.Phone,
	}

	_, err = svc.userRepo.Update(updated)
	if err != nil {
		svc.log.Error("cannot update user", "error", err.Error())
		return err
	}

	return nil
}

func (svc *userService) GetUser(publicID string) (*gormodel.User, error) {
	user, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "public_id", publicID)
		return nil, errors.New("user not found")
	}
	return user, nil
}
