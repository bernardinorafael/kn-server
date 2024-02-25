package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID
	Name       string
	Username   string
	Email      string
	PersonalID string
	Password   string
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
