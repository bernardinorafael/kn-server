package entity

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Username   string         `json:"username" gorm:"unique"`
	Email      string         `json:"email" gorm:"unique"`
	Password   string         `json:"password,omitempty"`
	PersonalID string         `json:"personal_id" gorm:"unique"`
	Active     bool           `json:"active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}
