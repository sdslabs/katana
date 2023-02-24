package controllers

import (
	"fmt"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func ChallengeUpdate(c *fiber.Ctx) error {

	cmd := exec.Command("/bin/sh", "./teams/team1/script.sh")

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))

	type Teams struct {
		URL    string `json:"url"`
		BEFORE string `json:"before"`
	}
	var req Teams
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse JSON request body: %v", err))
	}

	// Print out the ID value
	fmt.Printf("Received : %d\n", req)

	//return c.SendString(fmt.Sprintf("Received full_name: %s", req.BEFORE))
	return c.SendString("Working")
}
