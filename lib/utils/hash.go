package utils

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
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
	return err == nil
}

// EncodePassword encodes password using PBKDF2 SHA256 with given salt.
func EncodePassword(password, salt string) string {
	newPasswd := pbkdf2.Key([]byte(password), []byte(salt), 10000, 50, sha256.New)
	return fmt.Sprintf("%x", newPasswd)
}
