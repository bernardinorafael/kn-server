package service

import (
	"fmt"
	"log/slog"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
)

type userService struct {
	log      *slog.Logger
	userRepo contract.UserRepository
}

func NewUserService(log *slog.Logger, userRepo contract.UserRepository) contract.UserService {
	return &userService{log: log, userRepo: userRepo}
}

func (svc *userService) GetUser(publicID string) (*model.User, error) {
	u, err := svc.userRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error(fmt.Sprintf("product with PublicID [%s] not found", publicID))
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}
