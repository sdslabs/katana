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
		fmt.Errorf("Error building docker image: %s", err.Error())
	}

	// Push the docker image to Harbor
	cmd = exec.Command("docker", "push", configs.HarborConfig.Hostname+"/katana/"+folderName)
	if err := cmd.Run(); err != nil {
		fmt.Errorf("Error pushing docker image: %s", err.Error())
	}
}

func DockerLogin(username string, password string) error {
	cmd := fmt.Sprintf("sudo docker login -u %s -p %s %s", username, password, configs.KatanaConfig.Harbor.Hostname)
	out := exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	return nil
}
