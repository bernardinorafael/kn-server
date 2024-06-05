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

type UpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
