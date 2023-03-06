package controllers

import (
	"fmt"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func ChallengeUpdate(c *fiber.Ctx) error {

	type Teams struct {
		TEAMNAME string `json:"url"`
		BEFORE   string `json:"before"`
	}
	var req Teams
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse JSON request body: %v", err))
	}

	teamname := req.TEAMNAME

	// Print out the ID value
	fmt.Printf("Received : %v\n", req)
	fmt.Printf("Test : %s\n", teamname)

	cmd := exec.Command("/bin/sh", "./teams/team1/script.sh")

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))

	//return c.SendString(fmt.Sprintf("Received full_name: %s", req.BEFORE))
	return c.SendString("Working")
}
