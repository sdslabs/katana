package controllers

import (
	"fmt"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func ChallenegUpdate(c *fiber.Ctx) error {

	cmd := exec.Command("/bin/sh", "../../../teams/team1/script.sh")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	// Print the output
	fmt.Println(string(stdout))

	return c.SendString("Challenge successfully pulled from gi")
}
