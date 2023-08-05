package utils

import (
	"log"

	"github.com/sdslabs/katana/configs"
)

func BuildDockerImage(_ChallengeName string, _DockerfilePath string) {
	log.Println("Building docker image, Please wait...")
	cmd := "docker build -t " + configs.HarborConfig.Hostname + "/katana/" + _ChallengeName + " " + _DockerfilePath
	if err := RunCommand(cmd); err != nil {
		log.Printf("Error: %s\n", err)
	}
	log.Println("Docker image built successfully")

	log.Println("Pushing docker image, Please wait...")
	cmd = "docker push " + configs.HarborConfig.Hostname + "/katana/" + _ChallengeName
	if err := RunCommand(cmd); err != nil {
		log.Printf("Error: %s\n", err)
	}
	log.Println("Docker image pushed successfully")
}

func DockerLogin(username string, password string) {
	log.Println("Logging into Harbor, Please wait...")
	cmd := "docker login -u " + username + " -p " + password + " " + configs.KatanaConfig.Harbor.Hostname
	if err := RunCommand(cmd); err != nil {
		log.Printf("Error: %s\n", err)
	}
	log.Println("Logged into Harbor successfully")
}
