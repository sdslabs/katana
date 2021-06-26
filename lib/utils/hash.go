package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pass := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashWithPassword(hashedPassword, password string) bool {
	hash := []byte(hashedPassword)
	pass := []byte(password)
	err := bcrypt.CompareHashAndPassword(hash, pass)
	if err != nil {
		return false
	}
	return true
}
