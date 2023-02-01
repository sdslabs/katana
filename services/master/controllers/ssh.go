package controllers

import (
	"github.com/gofiber/fiber/v2"
	ssh "github.com/sdslabs/katana/services/sshproviderservice"
)

func SSH(c *fiber.Ctx) error {
	go func() {
		x := ssh.Server()
		x.ListenAndServe()
	}()
	ssh.CreateTeams()
	return c.SendString("ok")
}
