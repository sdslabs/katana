package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/sdslabs/katana/configs"
)

func BuildDockerImage(_ChallengeName string, _DockerfilePath string) {
	buf := new(bytes.Buffer)
	if err := Tar(_DockerfilePath, buf); err != nil {
		log.Fatal(err, ": error tarring directory")
		return
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Building Docker image, Please wait......")
	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		buf,
		types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Remove:     true,
			Tags:       []string{"harbor.katana.local/katana/" + _ChallengeName},
		},
	)
	if err != nil {
		log.Fatal(err, " :unable to create image")
		return
	}

	defer imageBuildResponse.Body.Close()
	log.Println("Docker image built successfully")
	DockerLogin(configs.KatanaConfig.Harbor.Username, configs.KatanaConfig.Harbor.Password)
	log.Println("Pushing Docker image to Harbor, please wait...")

	authConfig := registry.AuthConfig{
		Username: configs.KatanaConfig.Harbor.Username,
		Password: configs.KatanaConfig.Harbor.Password,
	}
	authJSON, err := json.Marshal(authConfig)
	if err != nil {
		log.Fatal(err, ": error encoding credentials")
		return
	}

	encodedAuth := Base64Encode(string(authJSON))

	pushOptions := types.ImagePushOptions{
		RegistryAuth: encodedAuth,
	}

	pushResponse, err := cli.ImagePush(ctx, "harbor.katana.local/katana/"+_ChallengeName, pushOptions)
	if err != nil {
		log.Fatal(err, " :unable to push docker image")
		return
	}
	defer pushResponse.Close()

}

func DockerLogin(username string, password string) {
	log.Println("Logging into Harbor, Please wait...")

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error creating Docker client:", err)
		return
	}

	authConfig := registry.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: "harbor.katana.local",
	}

	_, err = cli.RegistryLogin(ctx, authConfig)
	if err != nil {
		log.Printf("Error during login: %s\n", err)
		return
	}
	if err != nil {
		log.Fatal(err, " :unable to read push response")
		return
	}

	log.Println("Logged into Harbor successfully")
}

func DockerImageExists(imageName string) bool {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error: ", err)
		return false
	}

	_, _, err = cli.ImageInspectWithRaw(ctx, "harbor.katana.local/katana/"+imageName)
	return err == nil
}
