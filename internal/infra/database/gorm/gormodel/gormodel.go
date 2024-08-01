package gormodel

import (
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/cpf"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/password"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/phone"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/slug"

	"gorm.io/gorm"
)

type Team struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	PublicID  string         `json:"public_id" gorm:"unique"`
	Name      string         `json:"name"`
	Members   []User         `json:"members" gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type User struct {
	ID        int               `json:"id" gorm:"primaryKey"`
	PublicID  string            `json:"public_id" gorm:"unique"`
	Name      string            `json:"name"`
	Email     email.Email       `json:"email" gorm:"unique"`
	Document  cpf.CPF           `json:"document" gorm:"unique"`
	Phone     phone.Phone       `json:"phone" gorm:"unique"`
	Enabled   bool              `json:"enabled"`
	Password  password.Password `json:"password"`
	TeamID    *int              `json:"team_id" gorm:"default:null"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `json:"deleted_at" gorm:"index"`
}

type Product struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	PublicID  string         `json:"public_id" gorm:"unique"`
	Slug      slug.Slug      `json:"slug" gorm:"unique"`
	Name      string         `json:"name" gorm:"index"`
	Image     string         `json:"image"`
	Price     float64        `json:"price"`
	Quantity  int            `json:"quantity"`
	Enabled   bool           `json:"enabled gorm:index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
