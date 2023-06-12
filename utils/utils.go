package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}

func VerifyPassword(loginPassword string, dbPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(loginPassword))
}
