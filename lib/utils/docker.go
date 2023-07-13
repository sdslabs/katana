package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sdslabs/katana/configs"
)

func BuildDockerImage(_DockerfilePath string) {
	// Remove Dockerfile from the path
	folderPath := _DockerfilePath[:len(_DockerfilePath)-11]
	folderName := folderPath[strings.LastIndex(folderPath, "/")+1:]

	// Build the docker image
	cmd := exec.Command("docker", "build", "-t", configs.HarborConfig.Hostname+"/katana/"+folderName, folderPath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	// Push the docker image to Harbor
	cmd = exec.Command("docker", "push", configs.HarborConfig.Hostname+"/katana/"+folderName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func DockerLogin(username string, password string) {
	cmd := fmt.Sprintf("docker login -u %s -p %s %s", username, password, configs.KatanaConfig.Harbor.Hostname)
	out := exec.Command("sh", "-c", cmd)
	if err := out.Run(); err != nil {
		fmt.Printf("Docker login error: %s\n", err)
	}
}
