package contract

import (
	"context"

	"github.com/bernardinorafael/gozinho/internal/application/dto"
	"github.com/bernardinorafael/gozinho/internal/infra/rest/response"
)

type AccountService interface {
	Save(ctx context.Context, u *dto.UserInput) error
	GetByID(ctx context.Context, id string) (*response.UserResponse, error)
	Update(ctx context.Context, u *dto.UpdateUser, id string) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) (*response.AllUsersResponse, error)
	UpdatePassword(ctx context.Context, u *dto.UpdatePassword, id string) error
}
