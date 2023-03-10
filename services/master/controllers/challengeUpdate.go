package controllers

import (
	"context"
	"fmt"
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
	challengeName := s[1]
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

	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		fmt.Println(err)
	}

	// Build the challenge image with the Dockerfile in the challenge directory
	out, err := cli.ImageBuild(context.Background(), strings.NewReader("FROM scratch"), types.ImageBuildOptions{
		Dockerfile: dir + "/Dockerfile",
		Tags:       []string{dir},
	})
	if err != nil {
		fmt.Println(err)
	}

	defer out.Body.Close()

	// Create a labelSelector to get the challenge pod
	labelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": teamName + "_" + challengeName,
		},
	}

	// Delete the challenge pod
	err = client.CoreV1().Pods(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&labelSelector),
	})
	if err != nil {
		fmt.Println(err)
	}

	return c.SendString("Challenge updated")

}
