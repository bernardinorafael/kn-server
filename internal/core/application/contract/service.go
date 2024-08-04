package contract

import (
	"io"

	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
)

type AuthService interface {
	Login(dto dto.Login) (gormodel.User, error)
	LoginOTP(dto dto.LoginOTP) (gormodel.User, error)
	Register(dto dto.Register) (gormodel.User, error)
}

type ProductService interface {
	Create(dto dto.CreateProduct) error
	Delete(publicID string) error
	GetByPublicID(publicId string) (gormodel.Product, error)
	GetBySlug(slugInput string) (gormodel.Product, error)
	GetAll(disabled bool, orderBy string) ([]gormodel.Product, error)
	UpdatePrice(publicID string, price int) error
	IncreaseQuantity(publicID string, quantity int) error
	ChangeStatus(publicID string, status bool) error
}

type UserService interface {
	GetUser(publicID string) (gormodel.User, error)
	Update(publicID string, dto dto.UpdateUser) error
	RecoverPassword(publicID string, dto dto.UpdatePassword) error
}

type TeamService interface {
	Create(dto dto.CreateTeam) error
	GetByID(publicID string) (gormodel.Team, error)
}

type FileManagerService interface {
	UploadFile(file io.Reader, key, bucket string) (location string, err error)
}

type SMSNotifier interface {
	Notify(to string) error
	Confirm(code string, phone string) error
}
