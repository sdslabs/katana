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
	cmd := exec.Command("docker", "build", "-t", configs.HarborConfig.Hostname+"/katana/"+_ChallengeName, _DockerfilePath)
	fmt.Println("docker build -t", configs.HarborConfig.Hostname+"/katana/"+_ChallengeName, _DockerfilePath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Println("Docker image built successfully")

	// Push the docker image to Harbor
	cmd = exec.Command("docker", "push", configs.HarborConfig.Hostname+"/katana/"+_ChallengeName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Println("Docker image pushed successfully")
}

func DockerLogin(username string, password string) {
	cmd := fmt.Sprintf("docker login -u %s -p %s %s", username, password, configs.KatanaConfig.Harbor.Hostname)
	out := exec.Command("sh", "-c", cmd)
	if err := out.Run(); err != nil {
		fmt.Printf("Docker login error: %s\n", err)
	}
}
