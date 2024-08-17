package gormodel

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	PublicID  string         `json:"public_id" gorm:"unique"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"unique"`
	Phone     string         `json:"phone" gorm:"unique"`
	Status    string         `json:"status"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Product struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	PublicID  string         `json:"public_id" gorm:"unique"`
	Slug      string         `json:"slug" gorm:"unique"`
	Name      string         `json:"name" gorm:"index"`
	Image     string         `json:"image"`
	Price     int            `json:"price"`
	Quantity  int            `json:"quantity"`
	Enabled   bool           `json:"enabled gorm:index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
