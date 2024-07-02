package response

import (
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/cpf"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/email"
	"github.com/bernardinorafael/kn-server/internal/core/domain/valueobj/slug"
)

type Product struct {
	PublicID  string    `json:"public_id"`
	Slug      slug.Slug `json:"slug"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int32     `json:"quantity"`
	Enabled   bool      `json:"enabled"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	Name      string      `json:"name"`
	Email     email.Email `json:"email"`
	PublicID  string      `json:"public_id"`
	Document  cpf.CPF     `json:"document"`
	Enabled   bool        `json:"enabled"`
	CreatedAt time.Time   `json:"created_at"`
}
