package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GetTokenIssuer retrieves the TOKEN_ISSUER environment variable.
func GetTokenIssuer() (string, error) {
	tokenIssuer := os.Getenv("TOKEN_ISSUER")
	if tokenIssuer == "" {
		return "", errors.New("TOKEN_ISSUER not set")
	}

	return tokenIssuer, nil
}

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

// MakeJWT creates a new JWT.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	tokenIssuer, _ := GetTokenIssuer()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    tokenIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %s", err)
	}

	return signedToken, nil
}

// validateJWT validates a JWT by extracting the claims and checking
// the issuer, expiration time, and user ID.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	tokenIssuer, _ := GetTokenIssuer()

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != tokenIssuer {
		return uuid.Nil, fmt.Errorf("invalid issuer: %s", issuer)
	}

	expirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return uuid.Nil, err
	}
	if expirationTime.Before(time.Now().UTC()) {
		return uuid.Nil, fmt.Errorf("JWT is expired")
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %s", err)
	}

	return id, nil
}

// MakeRefreshToken makes a random 256 bit token encoded in hex.
func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)

	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}
