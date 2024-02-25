package dto

type CreateUserDTO struct {
	Name       string `json:"name" validate:"required,min=3,max=30"`
	Username   string `json:"username" validate:"required,min=3,max=14"`
	Email      string `json:"email" validate:"required,email"`
	PersonalID string `json:"personal_id" validate:"required,len=11"`
	Password   string `json:"password" validate:"required,min=8"`
}
