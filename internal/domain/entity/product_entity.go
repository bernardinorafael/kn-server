package entity

import (
	"time"

	"gorm.io/gorm"
)

type color struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Hex   string `json:"hex"`
}

type size struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	Availability string `json:"availability"`
}

type Product struct {
	ID          uint    `json:"id" gorm:"autoIncrement"`
	Name        string  `json:"name" gorm:"unique"`
	Slug        string  `json:"slug" gorm:"unique"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	SKU         string  `json:"sku" gorm:"unique"`

	Keywords []string `json:"keywords"`

	Color color  `json:"color"`
	Sizes []size `json:"sizes"`

	CategoryID string `json:"category_id"`
	StoreID    string `json:"store_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
