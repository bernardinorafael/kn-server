package services

import (
	"github.com/bernardinorafael/gozinho/internal/application/interfaces"
	"github.com/bernardinorafael/gozinho/internal/domain/dto"
)

type UserService struct {
	repository interfaces.UserRepository
}

func NewUserService(repository interfaces.UserRepository) *UserService {
	return &UserService{repository}
}

func (s *UserService) Create(u *dto.CreateUserDTO) error {
	return nil
}
