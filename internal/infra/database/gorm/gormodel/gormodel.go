package gormodel

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	PublicID  string         `json:"public_id" gorm:"unique"`
	OwnerID   string         `json:"owner_id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Members []User `json:"members" gorm:"foreignKey:PublicTeamID"`
}

type User struct {
	ID           int            `json:"id" gorm:"primaryKey"`
	PublicID     string         `json:"public_id" gorm:"unique"`
	Name         string         `json:"name"`
	Email        string         `json:"email" gorm:"unique"`
	Document     string         `json:"document" gorm:"unique"`
	Phone        string         `json:"phone" gorm:"unique"`
	Enabled      bool           `json:"enabled"`
	Password     string         `json:"password"`
	PublicTeamID *string        `json:"public_team_id" gorm:"default:null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Product struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	PublicID  string         `json:"public_id" gorm:"unique"`
	Slug      string         `json:"slug" gorm:"unique"`
	Name      string         `json:"name" gorm:"index"`
	Image     string         `json:"image"`
	Price     float64        `json:"price"`
	Quantity  int            `json:"quantity"`
	Enabled   bool           `json:"enabled gorm:index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
