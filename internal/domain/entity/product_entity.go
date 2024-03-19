package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        string         `json:"id"`
	Name      string         `json:"name" gorm:"unique"`
	Price     float64        `json:"price"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
