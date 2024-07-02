package contract

import (
	"io"

	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
)

type AuthService interface {
	Login(mail, password string) (*model.User, error)
	Register(name, mail, password, document string) (*model.User, error)
	RecoverPassword(publicID string, data dto.UpdatePassword) error
}

type ProductService interface {
	Create(data dto.CreateProduct, file io.Reader, fileName string) error
	Delete(publicID string) error
	GetByPublicID(publicID string) (*model.Product, error)
	GetBySlug(slugInput string) (*model.Product, error)
	GetAll() ([]model.Product, error)
	UpdatePrice(publicID string, price float64) error
	IncreaseQuantity(publicID string, quantity int32) error
}

type UserService interface {
	GetUser(publicID string) (*model.User, error)
}
