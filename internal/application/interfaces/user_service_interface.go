package interfaces

import "github.com/bernardinorafael/gozinho/internal/domain/dto"

type UserService interface {
	Create(u *dto.CreateUserDTO) error
}
