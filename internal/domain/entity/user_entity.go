package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Surname   string         `json:"surname"`
	Password  string         `json:"password,omitempty"`
	Email     string         `json:"email" gorm:"unique"`
	Document  string         `json:"document" gorm:"unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
