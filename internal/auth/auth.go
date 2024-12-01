package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword uses bcrypt to generate a hashed password from plaintext.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("couldn't hash password: %s", err)
	}

	return string(hash), err
}

// CheckPasswordHash compares the plaintext password with the hashed password.
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("password doesn't match hash")
	}

	return nil
}
