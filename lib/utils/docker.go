package utils

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

func createTarball(sourceDir string, buf *bytes.Buffer) error {
	tw := tar.NewWriter(buf)
	defer tw.Close()

	err := filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tw, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func BuildDockerImage(_ChallengeName string, _DockerfilePath string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	buildContextDir := _DockerfilePath

	err = createTarball(buildContextDir, buf) // Using the createTarball function
	if err != nil {
		log.Fatal(err, " :unable to create tarball")
	}

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		buf,
		types.ImageBuildOptions{
			Context:    buf,
			Dockerfile: "Dockerfile",
			Remove:     true,
			Tags:       []string{"images"},
		},
	)
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer imageBuildResponse.Body.Close()

	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
	// DockerLogin("username", "pass")
	// log.Println("Pushing Docker image to Docker Hub, please wait...")

	// authConfig := registry.AuthConfig{
	// 	Username: "username",
	// 	Password: "pass",
	// }
	// encodedAuth, err := encodeAuthToBase64(authConfig)
	// if err != nil {
	// 	log.Printf("Error encoding auth config: %s\n", err)
	// 	return
	// }

	// pushOptions := types.ImagePushOptions{
	// 	RegistryAuth: encodedAuth,
	// }

	// pushResponse, err := cli.ImagePush(ctx, "harbor.katana.local/katana/"+_ChallengeName, pushOptions)
	// if err != nil {
	// 	log.Fatal(err, " :unable to push docker image")
	// }
	// defer pushResponse.Close()

	// _, err = io.Copy(os.Stdout, pushResponse)
	// if err != nil {
	// 	log.Fatal(err, " :unable to read push response")
	// }
}

// func encodeAuthToBase64(authConfig registry.AuthConfig) (string, error) {
// 	authJSON, err := json.Marshal(authConfig)
// 	if err != nil {
// 		return "", err
// 	}
// 	return base64.URLEncoding.EncodeToString(authJSON), nil
// }

func DockerLogin(username string, password string) {
	log.Println("Logging into Harbor, Please wait...")

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Error creating Docker client: %s\n", err)
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
	// _, err = io.Copy(os.Stdout, loginResponse)
	if err != nil {
		log.Fatal(err, " :unable to read push response")
	}

	log.Println("Logged into Harbor successfully")
}

func DockerImageExists(imageName string) bool {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Error: %s\n", err)
		return false
	}

	_, _, err = cli.ImageInspectWithRaw(ctx, "harbor.katana.local/katana/"+imageName)
	return err == nil
}
