package dto

import "io"

type CreateProduct struct {
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Quantity  int       `json:"quantity"`
	Image     io.Reader `json:"image"`
	ImageName string    `json:"image_name"`
}

type UpdatePrice struct {
	Amount int `json:"amount"`
}

type UpdateQuantity struct {
	Amount int `json:"amount"`
}

type ChangeStatus struct {
	Status bool `json:"status"`
}
