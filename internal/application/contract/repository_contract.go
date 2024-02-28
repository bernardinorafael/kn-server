package contract

import (
	"github.com/bernardinorafael/gozinho/internal/domain/entity"
)

type AccountRepository interface {
	Save(u *entity.Account) error
	GetByEmail(email string) (*entity.Account, error)
	GetByID(id string) (*entity.Account, error)
	Update(u *entity.Account) error
	Delete(id string) error
	GetAll() ([]entity.Account, error)
	UpdatePassword(password string, id string) error
	CheckUserExist(email, username, personalID string) (bool, error)
	GetPassword(id string) (string, error)
}
