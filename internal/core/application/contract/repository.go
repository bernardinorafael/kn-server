package contract

import (
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
)

type GormUserRepository interface {
	Create(u gormodel.User) error
	GetByPublicID(publicID string) (gormodel.User, error)
	GetByEmail(email string) (gormodel.User, error)
	GetByPhone(phone string) (gormodel.User, error)
	Update(u gormodel.User) (gormodel.User, error)
}

type GormProductRepository interface {
	Create(p gormodel.Product) error
	GetByID(id int) (gormodel.Product, error)
	GetByPublicID(publicID string) (gormodel.Product, error)
	GetBySlug(name string) (gormodel.Product, error)
	GetAll(dto dto.ProductsFilter) ([]gormodel.Product, error)
	Delete(publicID string) error
	Update(p gormodel.Product) (gormodel.Product, error)
}
