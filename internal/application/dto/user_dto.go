package dto

type CreateUser struct {
	Name     string `json:"name" validate:"required,min=3,max=30"`
	Surname  string `json:"surname" validate:"required,min=5,max=5"`
	Email    string `json:"email" validate:"required,email"`
	Document string `json:"document" validate:"required,len=11,numeric"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUser struct {
	Name    string `json:"name" validate:"omitempty,min=3,max=30"`
	Surname string `json:"surname" validate:"required,min=2,max=5"`
	Email   string `json:"email" validate:"omitempty,email"`
}
