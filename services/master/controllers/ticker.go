package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdslabs/katana/lib/retrieve_data"
	t "github.com/sdslabs/katana/services/flaghandlerservice"
)

func StartTicker(c *fiber.Ctx) error {
	t.Server()
	retrieve_data.StartSaving()
	return c.SendString("Ticker Started")
}
