package controllers

import (
	"github.com/gofiber/fiber/v2"
	deployer "github.com/sdslabs/katana/services/challengedeployerservice"
)

func Deploy(c *fiber.Ctx) error {

	localFilePath := "/home/percy/notekeeper.tar.gz"
	pathInPod := "/opt/katana/katana_web_notekeeper.tar.gz"
	deployer.DeployToAll(localFilePath, pathInPod)

	return c.SendString("Deployed")
}
