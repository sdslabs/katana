package controllers

import (
	"github.com/gofiber/fiber/v2"
	deployer "github.com/sdslabs/katana/services/challengedeployerservice"
)

func Deploy(c *fiber.Ctx) error {

	localFilePath := "/home/paradox/katana/LICENSE"
	pathInPod := "/opt/katana/LICENSE"

	deployer.DeployToAll(localFilePath, pathInPod)

	return c.SendString("Deployed")
}
