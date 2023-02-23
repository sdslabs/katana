package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ChallengeUpdate(c *fiber.Ctx) error {

	//cmd := exec.Command("/bin/sh", "../../../teams/team1/script.sh")

	type Teams struct {
		URL  int    `json:"url"`
		Name string `json:"name"`
	}
	var req Teams
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse JSON request body: %v", err))
	}

	// Print out the ID value
	fmt.Printf("Received : %d\n", req.URL)

	return c.SendString(fmt.Sprintf("Received full_name: %s", req.Name))

}
