package memory

import (
	"errors"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"github.com/bernardinorafael/kn-server/internal/domain/repository"
	"github.com/google/uuid"
)

type userMemoryRepository struct {
	data map[string]*entity.User
}

func newUserMemoryRepository() repository.UserRepository {
	return &userMemoryRepository{data: make(map[string]*entity.User)}
}

func (u userMemoryRepository) Create(user *entity.User) error {
	user, err := entity.NewUser("john doe", "john_doe@gmail.com", "abcd1234")
	if err != nil {
		return err
	}
	u.data[user.ID.String()] = user

	return nil
}

func (u userMemoryRepository) FindByID(id uuid.UUID) (*entity.User, error) {
	user, ok := u.data[id.String()]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}
