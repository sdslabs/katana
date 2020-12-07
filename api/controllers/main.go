package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func HelloMain(c *fiber.Ctx) error {
	msg := fmt.Sprintf("Hello, admin %s 👋!", c.Params("name"))
	return c.SendString(msg)
}
