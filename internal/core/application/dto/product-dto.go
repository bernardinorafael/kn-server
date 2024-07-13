package dto

import "io"

type CreateProduct struct {
	Name      string
	Price     float64
	Quantity  int32
	Image     io.Reader
	ImageName string
}

type UpdatePrice struct {
	Amount float64
}

type UpdateQuantity struct {
	Amount int32
}
