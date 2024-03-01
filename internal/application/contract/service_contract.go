package contract

import (
	"context"

	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/response"
)

type AccountService interface {
	Save(ctx context.Context, u dto.UserInput) error
	GetByID(ctx context.Context, id string) (*response.UserResponse, error)
	Update(ctx context.Context, u dto.UpdateUser, id string) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) (*response.AllUsersResponse, error)
	UpdatePassword(ctx context.Context, u dto.UpdatePassword, id string) error
	Login(ctx context.Context, u dto.Login) (*entity.Account, error)
}
