package crypto

import "golang.org/x/crypto/bcrypt"

func Make(password string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}

func Compare(encrypted, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
