package entity

import (
	"github.com/google/uuid"
	"time"

	"github.com/bernardinorafael/kn-server/helper/crypto"
	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email" gorm:"unique"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func NewUser(name, pass, email string) (*User, error) {
	encrypted, err := crypto.Make(pass)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.New().String(),
		Email:    email,
		Name:     name,
		Password: encrypted,
	}, nil
}
