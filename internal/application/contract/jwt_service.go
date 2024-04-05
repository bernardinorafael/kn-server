package contract

import (
	"github.com/bernardinorafael/kn-server/internal/application/dto"
)

type JWTService interface {
	CreateToken(id string) (string, *dto.Claims, error)
	ValidateToken(token string) (*dto.Claims, error)
}
