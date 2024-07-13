package service

import (
	"fmt"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
	"github.com/bernardinorafael/kn-server/pkg/logger"
)

type userService struct {
	log      logger.Logger
	userRepo contract.UserRepository
}

func NewUserService(log logger.Logger, userRepo contract.UserRepository) contract.UserService {
	return &userService{log: log, userRepo: userRepo}
}

func (svc *userService) GetUser(publicID string) (*model.User, error) {
	u, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("user not found", "id", publicID)
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}
