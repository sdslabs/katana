package controllers

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

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
	//http connection configuration for 30 min

	var p GogsRequest
	if err := c.BodyParser(&p); err != nil {
		return err
	}

	log.Println(p)
	dir := p.Repository.FullName
	log.Println(dir + " updated")
	s := strings.Split(dir, "/")
	challengeName := s[1]
	teamName := s[0]
	namespace := teamName + "-ns"

	repo, err := git.PlainOpen("teams/" + dir)
	if err != nil {
		fmt.Println(err)
	}

	auth := &http.BasicAuth{
		Username: g.AdminConfig.Username,
		Password: g.AdminConfig.Password,
	}

	worktree, err := repo.Worktree()
	resp := worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	log.Println(resp)
	if err != nil {
		fmt.Println("Error pulling changes:", err)
	}

	// Build the challenge with Dockerfile
	cmd := exec.Command("docker", "build", "-t", dir, "./teams/"+dir)
	cmd.Run()

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
