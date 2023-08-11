package utils

import (
	r "crypto/rand"
	"math/big"
	"math/rand"
)

func RandomString(n uint) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// RandomChars returns a generated string in given number of random characters.
func randomChars(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	randomInt := func(max *big.Int) (int, error) {
		r, err := r.Int(r.Reader, max)
		if err != nil {
			return 0, err
		}

		return int(r.Int64()), nil
	}

	buffer := make([]byte, n)
	max := big.NewInt(int64(len(alphanum)))
	for i := 0; i < n; i++ {
		index, err := randomInt(max)
		if err != nil {
			return "", err
		}

		buffer[i] = alphanum[index]
	}

	return string(buffer), nil
}

// RandomSalt returns randomly generated 10-character string that can be used as
// the user salt.
func RandomSalt() (string, error) {
	return randomChars(10)
}
