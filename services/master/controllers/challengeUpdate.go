package controllers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gofiber/fiber/v2"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Repo struct {
	FullName string `json:"full_name"`
}

type GogsRequest struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	Repository Repo   `json:"repository"`
}

//	func Chall() string {
//		cmd, err := exec.Command("/bin/sh", "asdf.sh").Output()
//		if err != nil {
//			fmt.Printf("error %s", err)
//		}
//		output := string(cmd)
//		return output
//	}
func ChallengeUpdate(c *fiber.Ctx) error {

	client, err := utils.GetKubeClient()
	if err != nil {
		fmt.Println(err)
	}

	var p GogsRequest
	if err := c.BodyParser(&p); err != nil {
		return err
	}

	dir := p.Repository.FullName
	s := strings.Split(dir, "/")
	name := s[1]
	teamName := s[0]
	namespace := teamName + "-ns"

	repo, err := git.PlainOpen(dir)
	if err != nil {
		fmt.Println(err)
	}

	auth := &http.BasicAuth{
		Username: g.AdminConfig.Username,
		Password: g.AdminConfig.Password,
	}

	worktree, err := repo.Worktree()
	worktree.Pull(&git.PullOptions{
		RemoteName: "origin master",
		Auth:       auth,
	})

	if err != nil {
		fmt.Println("Error pulling changes:", err)
	}

	// Build the challenge with Dockerfile

	cmd := exec.Command("docker", "build", "-t", dir, dir)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}
	dockerfile, err := os.Open(fmt.Sprintf("%s/Dockerfile", dir))
	if err != nil {
		panic(err)
	}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerfile, types.ImageBuildOptions{
		Tags:        []string{dir},
		Context:     dockerfile, // Use the Dockerfile as the build context
		Remove:      true,       // Remove intermediate containers after a successful build
		ForceRemove: true,       // Remove the image if it already exists
	})
	if err != nil {
		panic(err)
	}
	defer buildResponse.Body.Close()
	// Restart the pod using go client
	err = client.CoreV1().Pods(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println(err)
	}

	return c.SendString("Challenge updated")

}
