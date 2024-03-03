package contract

import (
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type AccountRepository interface {
	Save(u entity.Account) error
	GetByEmail(email string) (*entity.Account, error)
	GetByID(id string) (*entity.Account, error)
	Update(account *entity.Account, id string) error
	Delete(id string) error
	GetAll() (*[]entity.Account, error)
	UpdatePassword(password string, id string) error
	GetPassword(id string) (string, error)
}
