package main

import (
	"fmt"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func git(arg string) {
	app := "git"
	arg0 := arg
	cmd := exec.Command(app, arg0)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

func ChallengeUpdate(c *fiber.Ctx) error {
	git("pull")
	//TODO : Apply yaml file and up challenge in container in docker
}
