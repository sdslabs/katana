package controllers

import (
	"fmt"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

type Repo struct {
	FullName string `json:"full_name"`
}

type GogsRequest struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	Repository Repo   `json:"repository"`
}

func ChallengeUpdate(c *fiber.Ctx) error {

	var p GogsRequest
	if err := c.BodyParser(&p); err != nil {
		return err
	}

	dir := p.Repository.FullName

	// Git pull in the directory
	cmd := exec.Command("git", "-C", dir, "pull")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// Build the challenge with Dockerfile
	cmd = exec.Command("docker", "build", "-t", dir, dir)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// Delete the old manifest
	cmd = exec.Command("kubectl", "delete", "-f", dir+"/challenge.yaml")
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// Apply the new manifest
	cmd = exec.Command("kubectl", "apply", "-f", dir+"/challenge.yaml")
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	return c.SendString("Challenge updated")
}
