package dto

type Login struct {
	Email    string
	Password string
}

type Register struct {
	Name     string
	Email    string
	Password string
}

type CreateProduct struct {
	Name     string
	Price    float64
	Quantity int32
}
