package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// Logs returns the logs of the cluster
func Logs(c *fiber.Ctx) error {
	return c.SendString("Logs")
}
