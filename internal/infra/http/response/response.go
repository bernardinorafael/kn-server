package response

import (
	"time"
)

type Product struct {
	PublicID  string    `json:"public_id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Quantity  int       `json:"quantity"`
	Enabled   bool      `json:"enabled"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	PublicID  string    `json:"public_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Document  string    `json:"document"`
	Phone     string    `json:"phone"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

type Team struct {
	PublicID  string    `json:"public_id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	Members   []User    `json:"members"`
	CreatedAt time.Time `json:"created_at"`
}
