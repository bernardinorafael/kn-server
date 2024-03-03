package dto

type CreateAccount struct {
	Name     string `json:"name" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Document string `json:"document" validate:"required,len=11"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateAccount struct {
	Name  string `json:"name" validate:"omitempty,min=3,max=30"`
	Email string `json:"email" validate:"omitempty,email"`
}

type UpdatePassword struct {
	Password    string `json:"password" validate:"required,min=8"`
	OldPassword string `json:"old_password" validate:"eqfield=Password"`
}
