package dto

type CreateUser struct {
	Name     string `json:"name" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Document int    `json:"document" validate:"required,len=11,gte=0"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUser struct {
	Name  string `json:"name" validate:"omitempty,min=3,max=30"`
	Email string `json:"email" validate:"omitempty,email"`
}
