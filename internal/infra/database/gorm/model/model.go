package model

import (
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/cpf"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/password"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/slug"

	"gorm.io/gorm"
)

type User struct {
	ID        int               `json:"id" gorm:"primaryKey"`
	PublicID  string            `json:"public_id" gorm:"unique"`
	Name      string            `json:"name"`
	Email     email.Email       `json:"email" gorm:"unique"`
	Document  cpf.CPF           `json:"document" gorm:"unique"`
	Enabled   bool              `json:"enabled"`
	Password  password.Password `json:"password"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt gorm.DeletedAt    `json:"deleted_at"`
}

type Product struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	PublicID  string         `json:"public_id" gorm:"unique"`
	Slug      slug.Slug      `json:"slug" gorm:"unique"`
	Name      string         `json:"name"`
	Image     string         `json:"image"`
	Price     float64        `json:"price"`
	Quantity  int32          `json:"quantity"`
	Enabled   bool           `json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
