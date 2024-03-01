package response

import (
	"time"
)

type UserResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	PersonalID string    `json:"personal_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type AllUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type LoginResponse struct {
	UserID      string    `json:"user_id"`
	AccessToken string    `json:"access_token"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}
