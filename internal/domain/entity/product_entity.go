package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `json:"id" gorm:"autoIncrement"`
	Name        string         `json:"name" gorm:"unique"`
	Slug        string         `json:"slug" gorm:"unique"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Stock       int            `json:"stock"`
	SKU         string         `json:"sku" gorm:"unique"`
	Size        string         `json:"size"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}
