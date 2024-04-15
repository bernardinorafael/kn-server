package contract

import (
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type Auth interface {
	Register(i dto.Register) (*entity.User, error)
	Login(i dto.Login) (*entity.User, error)
}
