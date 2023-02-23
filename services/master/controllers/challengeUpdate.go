package controllers

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

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

	// Git pull in the directory
	cmd := exec.Command("git", "-C", dir, "pull")
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// Build the challenge with Dockerfile
	cmd = exec.Command("docker", "build", "-t", dir, dir)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// Restart the pod using go client
	err = client.CoreV1().Pods(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println(err)
	}

	return c.SendString("Challenge updated")
}
