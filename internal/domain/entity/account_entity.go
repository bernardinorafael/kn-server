package entity

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Username  string         `json:"username" gorm:"unique"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"password,omitempty"`
	Document  string         `json:"document" gorm:"unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
