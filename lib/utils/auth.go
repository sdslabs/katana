package utils

import (
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func VerifyToken(c *fiber.Ctx) bool {
	cookie := c.Cookies("jwt")

	if cookie == "" {
		return false
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	return true
}
