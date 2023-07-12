package controllers

import (
	"github.com/gofiber/fiber/v2"
	t"github.com/sdslabs/katana/services/flaghandlerservice"
)

func StartTicker(c *fiber.Ctx) error {
	t.Server()
	return c.SendString("Ticker Started")
}
