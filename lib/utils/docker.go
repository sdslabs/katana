package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"

	"github.com/sdslabs/katana/configs"
)

func dockerLogin(username string, password string) error {

	log.Println("Logging into Harbor, Please wait...")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error creating Docker client:", err)
		return err
	}

	authConfig := registry.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: "harbor.katana.local",
	}

	_, err = cli.RegistryLogin(context.Background(), authConfig)
	if err != nil {
		log.Printf("Error during login: %s\n", err)
		return err
	}
	log.Println("Logged into Harbor successfully")
	return nil
}

func CheckDockerfile(_DockerfilePath string) bool {
	_, err := os.Stat(_DockerfilePath + "/Dockerfile")
	return !os.IsNotExist(err)
}

func BuildDockerImage(_ChallengeName string, _DockerfilePath string)error {
	buf := new(bytes.Buffer)
	if err := Tar(_DockerfilePath, buf); err != nil {
		log.Fatal(err, ": error tarring directory")
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(_ChallengeName)
	log.Println("dfsdfsdfsdf")
	log.Println(_DockerfilePath)

	log.Println("Building Docker image, Please wait......")

	imageBuildResponse, err := cli.ImageBuild(
		context.Background(),
		buf,
		types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Remove:     true,
			Tags:       []string{"harbor.katana.local/katana/" + _ChallengeName},
		},
	)
	if err != nil {
		log.Fatal(err, " :unable to create image")
		return err
	}

	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
		return err
	}

	log.Println("Docker image built successfully")
	err = dockerLogin(configs.KatanaConfig.Harbor.Username, configs.KatanaConfig.Harbor.Password)
	if err != nil {
		return err
	}

	log.Println("Pushing Docker image to Harbor, please wait...")

	authConfig := registry.AuthConfig{
		Username: configs.KatanaConfig.Harbor.Username,
		Password: configs.KatanaConfig.Harbor.Password,
	}
	authJSON, err := json.Marshal(authConfig)
	if err != nil {
		log.Fatal(err, ": error encoding credentials")
		return err
	}

	encodedAuth := Base64Encode(string(authJSON))

	pushOptions := types.ImagePushOptions{
		RegistryAuth: encodedAuth,
	}

	_, err = cli.ImagePush(context.Background(), "harbor.katana.local/katana/"+_ChallengeName, pushOptions)
	if err != nil {
		log.Fatal(err, " :unable to push docker image")
		return err
	}

	log.Println("Image Pushed to Harbor successfully")
	return nil
}

func DockerImageExists(imageName string) bool {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error: ", err)
		return false
	}

	_, _, err = cli.ImageInspectWithRaw(context.Background(), "harbor.katana.local/katana/"+imageName)
	return err == nil
}
