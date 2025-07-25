package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return hash, nil
}

func ValidatePassword(hashed, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashed, plain)

	return err == nil
}
