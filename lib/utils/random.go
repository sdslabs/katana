package utils

import (
	"math/rand"

	"github.com/sdslabs/katana/configs"
)

func RandomString(n uint) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func GenPassword() string {
	return RandomString(configs.SSHProviderConfig.PasswordLen)
}
