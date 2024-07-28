package dto

type Login struct {
	Email    string
	Password string
}

type Register struct {
	Name     string
	Email    string
	Document string
	Phone    string
	Password string
}

type UpdatePassword struct {
	OldPassword string
	NewPassword string
}
