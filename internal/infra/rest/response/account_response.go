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
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

type AllUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type AuthToken struct {
	AccessToken string `json:"access_token"`
}
