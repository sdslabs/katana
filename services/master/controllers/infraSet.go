package controllers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func InfraSet(c *fiber.Ctx) error {

	// if !utils.VerifyToken(c) {
	// 	return c.SendString("Unauthorized")
	// }

	pathToCfg := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	fmt.Println(g.KatanaConfig.KubeNameSpace)

	config, err := clientcmd.BuildConfigFromFlags("", pathToCfg)
	if err != nil {
		log.Fatal(err)
	}

	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	if err = deployment.DeployCluster(config, kubeclient); err != nil {
		log.Fatal(err)
	}

	fmt.Println(kubeclient)

	return c.SendString("Infrastructure setup completed")
}
