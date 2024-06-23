package dto

type CreateProduct struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
}

type UpdatePrice struct {
	Amount float64 `json:"amount"`
}
