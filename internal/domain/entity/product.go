package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Slug     string  `json:"slug" gorm:"unique"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
	Status   bool    `json:"status" gorm:"default:true"`
}
