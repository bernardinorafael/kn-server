package interfaces

import (
	"github.com/bernardinorafael/gozinho/internal/application/dto"
	"github.com/bernardinorafael/gozinho/internal/infra/http/response"
	"github.com/google/uuid"
)

type UserService interface {
	Save(u *dto.CreateUser) error
	GetByID(id uuid.UUID) (*response.UserResponse, error)
	Update(u *dto.UpdateUser, id uuid.UUID) error
	Delete(id uuid.UUID) error
	GetAll() (*response.AllUsersResponse, error)
	UpdatePassword(u *dto.UpdatePassword, id uuid.UUID) error
}
