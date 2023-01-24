package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sdslabs/katana/lib/foundry"
	"github.com/sdslabs/katana/lib/utils"
)

func HelloAdmin(c *fiber.Ctx) error {

	if !utils.VerifyToken(c) {
		return c.SendStatus(403)
	}
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
