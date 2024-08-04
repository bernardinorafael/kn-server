package contract

import (
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/team"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
)

type UserRepository interface {
	Create(u user.User) (gormodel.User, error)
	GetByPublicID(publicID string) (gormodel.User, error)
	GetByEmail(email string) (gormodel.User, error)
	GetByPhone(phone string) (gormodel.User, error)
	Update(u user.User) (gormodel.User, error)
}

type ProductRepository interface {
	Create(p product.Product) (gormodel.Product, error)
	GetByID(id int) (gormodel.Product, error)
	GetByPublicID(publicID string) (gormodel.Product, error)
	GetBySlug(name string) (gormodel.Product, error)
	GetAll(disabled bool, orderBy string) ([]gormodel.Product, error)
	Delete(publicID string) error
	Update(p product.Product) (gormodel.Product, error)
}

type TeamRepository interface {
	Create(t team.Team) (gormodel.Team, error)
	Update(t team.Team) (gormodel.Team, error)
	Delete(publicID string) error
	GetByPublicID(publicID string) (gormodel.Team, error)
}
