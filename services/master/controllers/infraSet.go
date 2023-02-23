package controllers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func sshkeyExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func sshkeyCreate(path string) {
	//f, err := os.Create("testFile.txt")

	app := "ssh-keygen"
	arg0 := "-t"
	arg1 := "ed25519"
	arg2 := "-C"
	arg3 := "\"test@gmail.com\""
	arg4 := "-f"
	arg5 := path
	arg6 := "-q"
	arg7 := "-P"
	arg8 := "\"\""

	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))

	//defer f.Close()

}

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

	var filelocation string = "/tmp/ssh_admin"
	exist := sshkeyExists(filelocation)

	if exist {
		fmt.Println("Admin ssh file exists")
	} else {
		fmt.Println("Creating ssh keys")
		sshkeyCreate(filelocation)
	}

	return c.SendString("Infrastructure setup completed")
}
