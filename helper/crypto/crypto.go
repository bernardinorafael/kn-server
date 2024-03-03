package crypto

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}

func CheckPassword(encrypted, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
