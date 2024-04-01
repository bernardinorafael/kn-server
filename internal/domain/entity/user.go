package entity

import (
	"time"

	"github.com/bernardinorafael/kn-server/helper/crypto"
	"github.com/google/uuid"
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

func MakeUser(name, surname, pass, email, doc string) (*User, error) {
	encrypted, err := crypto.Make(pass)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       uuid.New().String(),
		Password: encrypted,
		Name:     name,
		Surname:  surname,
		Email:    email,
		Document: doc,
	}

	return user, nil
}
