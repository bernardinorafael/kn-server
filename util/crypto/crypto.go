package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(p), 12)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(encrypted), nil
}

func CheckPassword(p, encrypted string) error {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(p))
}
