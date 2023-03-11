package controllers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/gofiber/fiber/v2"
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
	//http connection configuration for 30 min

	var p GogsRequest
	if err := c.BodyParser(&p); err != nil {
		return err
	}

	log.Println(p)
	dir := "katana-team-0/notekeeper"
	log.Println(dir + " updated")
	s := strings.Split(dir, "/")
	challengeName := s[1]
	teamName := s[0]
	namespace := teamName + "-ns"

	repo, err := git.PlainOpen("teams/" + dir)
	if err != nil {
		fmt.Println(err)
	}

	// auth := &http.BasicAuth{
	// 	Username: g.AdminConfig.Username,
	// 	Password: g.AdminConfig.Password,
	// }

	worktree, err := repo.Worktree()
	resp := worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	log.Println(resp)
	if err != nil {
		fmt.Println("Error pulling changes:", err)
	}

	// Build the challenge with Dockerfile
	// cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation(), docker.WithHTTPClient(httpClient))
	// if err != nil {
	// 	log.Println(err)
	// }
	// buildCtx, err := os.Open("teams/" + dir + "/Dockerfile")
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer buildCtx.Close()
	// buildOptions := types.ImageBuildOptions{
	// 	Context:    buildCtx,
	// 	Dockerfile: "Dockerfile",
	// 	Tags:       []string{dir},
	// }
	// ctx := context.Background()
	// _, err = cli.ImageBuild(ctx, buildCtx, buildOptions)
	// if err != nil {
	// 	log.Println(err)
	// }
	// Create a labelSelector to get the challenge pod
	labelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": teamName + "_" + challengeName,
		},
	}
	log.Println("Yaha tak sexy 6")
	// Delete the challenge pod
	err = client.CoreV1().Pods(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&labelSelector),
	})
	if err != nil {
		log.Println(err)
		log.Println("Ye?")
		log.Println(err)
	}
	log.Println("Yaha tak sexy 7")
	return c.SendString("Challenge updated")

}
