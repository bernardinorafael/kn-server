package contract

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/model"
)

type AuthService interface {
	Login(data dto.Login) (*model.User, error)
	Register(data dto.Register) (*model.User, error)
	RecoverPassword(publicID string, data dto.UpdatePassword) error
}

type ProductService interface {
	Create(data dto.CreateProduct) error
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

// remove s3 deps from this contract
type FileManagerService interface {
	UploadFile(file io.Reader, key, bucket string) (*manager.UploadOutput, error)
}
