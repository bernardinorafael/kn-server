package dto

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOTP struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type Register struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type NotifySMS struct {
	Phone string `json:"phone"`
}
