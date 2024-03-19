package dto

type CreateProduct struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
	Size  string  `json:"size"`
}
