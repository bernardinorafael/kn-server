package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Password  string         `json:"password,omitempty"`
	Email     string         `json:"email" gorm:"unique"`
	Document  int            `json:"document" gorm:"unique"`
	Username  sql.NullString `json:"username" gorm:"unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
