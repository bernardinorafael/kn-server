package dto

type UserInput struct {
	Name     string `json:"name" validate:"required,min=3,max=30"`
	Username string `json:"username" validate:"required,min=3,max=14"`
	Email    string `json:"email" validate:"required,email"`
	Document string `json:"document" validate:"required,len=11"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUser struct {
	Name     string `json:"name" validate:"omitempty,min=3,max=30"`
	Username string `json:"username" validate:"omitempty,min=3,max=14"`
	Email    string `json:"email" validate:"omitempty,email"`
}

type UpdatePassword struct {
	Password    string `json:"password" validate:"required,min=8"`
	OldPassword string `json:"old_password" validate:"eqfield=Password"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
