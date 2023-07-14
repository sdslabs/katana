package utils

import (
	"fmt"

	"github.com/sdslabs/katana/configs"
)

func BuildDockerImage(_ChallengeName string, _DockerfilePath string) {
	// Build the docker image
	fmt.Println("Building docker image, Please wait...")
	cmd := "docker build -t " + configs.HarborConfig.Hostname + "/katana/" + _ChallengeName + " " + _DockerfilePath
	if err := RunCommand(cmd); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Println("Docker image built successfully")

	// Push the docker image to Harbor
	fmt.Println("Pushing docker image, Please wait...")
	cmd = "docker push " + configs.HarborConfig.Hostname + "/katana/" + _ChallengeName
	if err := RunCommand(cmd); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Println("Docker image pushed successfully")
}

func DockerLogin(username string, password string) {
	fmt.Println("Logging into Harbor, Please wait...")
	cmd := "docker login -u " + username + " -p " + password + " " + configs.KatanaConfig.Harbor.Hostname
	if err := RunCommand(cmd); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Println("Logged into Harbor successfully")
}
