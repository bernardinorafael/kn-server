package contract

import "github.com/bernardinorafael/kn-server/internal/domain/entity"

type AuthService interface {
	Login(email, password string) (*entity.User, error)
	Register(name, email, password string) (*entity.User, error)
}

type JWTService interface {
	CreateToken(id uint) (string, error)
	VerifyToken(token string) error
}
