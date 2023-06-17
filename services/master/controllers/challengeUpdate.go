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
	client, err := utils.GetClient("")
	if err != nil {
		fmt.Println(err)
	}
	//http connection configuration for 30 min

	var p GogsRequest
	if err := c.BodyParser(&p); err != nil {
		return err
	}

	dir := p.Repository.FullName
	s := strings.Split(dir, "/")
	challengeName := s[1]
	teamName := s[0]
	namespace := teamName + "-ns"
	log.Println("Challenge update request received for", challengeName, "by", teamName)
	repo, err := git.PlainOpen("teams/" + dir)
	if err != nil {
		fmt.Println(err)
	}

	auth := &http.BasicAuth{
		Username: g.AdminConfig.Username,
		Password: g.AdminConfig.Password,
	}

	worktree, err := repo.Worktree()
	worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	if err != nil {
		fmt.Println("Error pulling changes:", err)
	}

	log.Println("Pull successful for", teamName, ". Building image...")
	FirstPatch, err := exec.Command("docker", "inspect", dir).Output()
	cmd := exec.Command("docker", "build", "-t", dir, "./teams/"+dir)
	cmd.Run()
	cmd = exec.Command("minikube", "image", "load", dir)
	cmd.Run()
	if err != nil {
		log.Println(err)
	}
	if len(FirstPatch) <= 3 {
		log.Println("First Patch for", teamName)
		utils.DeployChallenge(challengeName, teamName, true)
	} else {
		log.Println("Not remaking", teamName)
		labelSelector := metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": challengeName,
			},
		}
		// Delete the challenge pod
		err = client.CoreV1().Pods(namespace).DeleteCollection(context.Background(), metav1.DeleteOptions{}, metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(&labelSelector),
		})
		if err != nil {
			log.Println("Error")
			log.Println(err)
		}
	}
	log.Println("Image built for", teamName)
	return c.SendString("Challenge updated")

}
