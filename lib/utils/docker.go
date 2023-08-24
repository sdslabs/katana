package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/sdslabs/katana/configs"
)

func Tar(src string, writers ...io.Writer) error {

	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix("./"+file, src+string(filepath.Separator))
		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})
}
func Untar(fileName string) io.Reader {
	tarGzFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer tarGzFile.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(tarGzFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer gzipReader.Close()

	// Read and store the content of the tar archive
	var contentBuf bytes.Buffer
	_, err = io.Copy(&contentBuf, gzipReader)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Return the contentBuf as an io.Reader
	return &contentBuf
}

func BuildDockerImage(_ChallengeName string, _DockerfilePath string) {
	buf := new(bytes.Buffer)
	if err := Tar(_DockerfilePath, buf); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(_DockerfilePath)
	tarName := _ChallengeName + ".tar.gz"
	f, err := os.Create(tarName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	if _, err := buf.WriteTo(f); err != nil {
		fmt.Println(err)
		return
	}
	reader := Untar(tarName)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		reader,
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

	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
	defer imageBuildResponse.Body.Close()

	DockerLogin(configs.KatanaConfig.Harbor.Username, configs.KatanaConfig.Harbor.Password)
	log.Println("Pushing Docker image to Docker Hub, please wait...")

	authConfig := registry.AuthConfig{
		Username: configs.KatanaConfig.Harbor.Username,
		Password: configs.KatanaConfig.Harbor.Password,
	}
	encodedAuth, err := encodeAuthToBase64(authConfig)
	if err != nil {
		log.Printf("Error encoding auth config: %s\n", err)
		return
	}

	pushOptions := types.ImagePushOptions{
		RegistryAuth: encodedAuth,
	}

	pushResponse, err := cli.ImagePush(ctx, "harbor.katana.local/katana/"+_ChallengeName, pushOptions)
	if err != nil {
		log.Fatal(err, " :unable to push docker image")
	}
	defer pushResponse.Close()

	_, err = io.Copy(os.Stdout, pushResponse)
	if err != nil {
		log.Fatal(err, " :unable to read push response")
	}
}

func encodeAuthToBase64(authConfig registry.AuthConfig) (string, error) {
	authJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(authJSON), nil
}

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
