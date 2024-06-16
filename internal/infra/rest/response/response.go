package response

import (
	"time"

	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/slug"
)

type Product struct {
	PublicID  string    `json:"public_id"`
	Slug      slug.Slug `json:"slug"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int32     `json:"quantity"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}
