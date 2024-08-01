package response

import (
	"time"
)

type Product struct {
	PublicID  string    `json:"public_id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Enabled   bool      `json:"enabled"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	PublicID  string    `json:"public_id"`
	Document  string    `json:"document"`
	Phone     string    `json:"phone"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}
