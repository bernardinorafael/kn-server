package contract

import (
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type AuthService interface {
	Register(input dto.Register) (*entity.User, error)
	Login(input dto.Login) (*entity.User, error)
}
