package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Username   string    `json:"username" gorm:"unique"`
	Email      string    `json:"email" gorm:"unique"`
	PersonalID string    `json:"personal_id" gorm:"unique"`
	Password   string    `json:"password,omitempty"`
	Active     bool      `json:"active" gorm:"default:true"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewUser(name, username, email, personalId, password string) *User {
	return &User{
		ID:         uuid.New(),
		Name:       name,
		Username:   username,
		Password:   password,
		Email:      email,
		PersonalID: personalId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Active:     true,
	}
}
