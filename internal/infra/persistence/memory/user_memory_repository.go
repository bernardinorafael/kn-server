package memory

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type userMemoryRepository struct {
	users []entity.User
}

func NewUserInMemoryRepository(db []entity.User) contract.UserRepository {
	return &userMemoryRepository{users: db}
}

func (repository *userMemoryRepository) Create(u entity.User) error {
	panic("method no implemented")
}

func (repository *userMemoryRepository) GetByEmail(email string) (*entity.User, error) {
	panic("method no implemented")
}

func (repository *userMemoryRepository) GetByID(id string) (*entity.User, error) {
	panic("method no implemented")
}

func (repository *userMemoryRepository) GetAll() (*[]entity.User, error) {
	panic("method no implemented")
}

func (repository *userMemoryRepository) GetPassword(id string) (string, error) {
	panic("method no implemented")
}
