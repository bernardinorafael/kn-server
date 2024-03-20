package authroute

import "time"

type LoginResponse struct {
	UserID      string    `json:"user_id"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	AccessToken string    `json:"access_token"`
}
