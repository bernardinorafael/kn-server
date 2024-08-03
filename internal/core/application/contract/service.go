package contract

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
)

type AuthService interface {
	Login(data dto.Login) (gormodel.User, error)
	Register(data dto.Register) (gormodel.User, error)
}

type ProductService interface {
	Create(data dto.CreateProduct) error
	Delete(publicID string) error
	GetByPublicID(publicId string) (gormodel.Product, error)
	GetBySlug(slugInput string) (gormodel.Product, error)
	GetAll(disabled bool, orderBy string) ([]gormodel.Product, error)
	UpdatePrice(publicID string, price float64) error
	IncreaseQuantity(publicID string, quantity int) error
	ChangeStatus(publicID string, status bool) error
}

type UserService interface {
	GetUser(publicID string) (gormodel.User, error)
	Update(publicID string, data dto.UpdateUser) error
	RecoverPassword(publicID string, data dto.UpdatePassword) error
}

type TeamService interface {
	Create(data dto.CreateTeam) error
	GetByID(publicID string) (gormodel.Team, error)
}

type FileManagerService interface {
	UploadFile(file io.Reader, key, bucket string) (*manager.UploadOutput, error)
}
