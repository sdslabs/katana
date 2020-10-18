package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sdslabs/katana/lib/foundry"
)

func HelloAdmin(c *fiber.Ctx) error {
	msg := fmt.Sprintf("Hello, admin %s 👋!", c.Params("name"))
	return c.SendString(msg)
}

func ClusterInfo(c *fiber.Ctx) error {
	response, err := foundry.ClusterInfo("abc", "localhost:8001")
	if err != nil {
		return err
	}
	return c.Send(response)
}
