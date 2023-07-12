package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	utils "github.com/sdslabs/katana/lib/utils"
)

func InfraSet(c *fiber.Ctx) error {

	// if !utils.VerifyToken(c) {
	// 	return c.SendString("Unauthorized")
	// }

	config, err := utils.GetKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	kubeclient, err := utils.GetKubeClient()
	if err != nil {
		log.Fatal(err)
	}

	// Generate the certificates
	generateCertsforHarbor()

	if err = deployment.DeployCluster(config, kubeclient); err != nil {
		log.Fatal(err)
	}

	fmt.Println(kubeclient)

	return c.SendString("Infrastructure setup completed")
}

func generateCertsforHarbor() {
	path, _ := os.Getwd()
	path = path + "/lib/harbor/certs"

	// Delete the directory if it already exists
	if _, err := os.Stat(path); os.IsExist(err) {
		errDir := os.RemoveAll(path)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		errDir := os.Mkdir(path, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	// Generate the certificates
	if err := utils.GenerateCerts(g.KatanaConfig.Harbor.Hostname, path); err != nil {
		log.Fatal(err)
	}
}
