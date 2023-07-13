package utils

import (
	"fmt"
	"os/exec"

	"github.com/sdslabs/katana/configs"
)

func BuildDockerImage(_ChallengeName string, _DockerfilePath string) {

	fmt.Println("folderpath :", _DockerfilePath)
	fmt.Println("foldername :", _ChallengeName)

	// Build the docker image
	fmt.Println("Building docker image, Please wait...")
	cmd := exec.Command("docker", "build", "-t", configs.HarborConfig.Hostname+"/katana/"+_ChallengeName, _DockerfilePath)
	fmt.Println("docker build -t", configs.HarborConfig.Hostname+"/katana/"+_ChallengeName, _DockerfilePath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Println("Docker image built successfully")

	// Push the docker image to Harbor
	fmt.Println("Pushing docker image, Please wait...")
	cmd = exec.Command("docker", "push", configs.HarborConfig.Hostname+"/katana/"+_ChallengeName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Println("Docker image pushed successfully")
}

func DockerLogin(username string, password string) {
	fmt.Println("Logging into Harbor, Please wait...")
	cmd := exec.Command("docker", "login", "-u", username, "-p", password, configs.KatanaConfig.Harbor.Hostname)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Docker login error: %s\n", err)
	}
	fmt.Println("Logged into Harbor successfully")
}
