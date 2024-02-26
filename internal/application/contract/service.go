package contract

import (
	"github.com/bernardinorafael/gozinho/internal/application/dto"
	"github.com/bernardinorafael/gozinho/internal/infra/http/response"
)

type UserService interface {
	Save(u *dto.UserInput) error
	GetByID(id string) (*response.UserResponse, error)
	Update(u *dto.UpdateUser, id string) error
	Delete(id string) error
	GetAll() (*response.AllUsersResponse, error)
	UpdatePassword(u *dto.UpdatePassword, id string) error
}
