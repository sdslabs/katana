package controllers

import (
	"fmt"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gofiber/fiber/v2"
	"github.com/sdslabs/katana/lib/utils"
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

	// fmt.Println(p.Repository)
	// fmt.Println(p.Repository.Name)
	fmt.Println("Request Recieved, please wait..pulling repo and building images")
	teamname := p.Repository.Name
	//fmt.Println(teamname)

	// Git pull in the directory
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

		cmd := exec.Command("/bin/sh", "./teams/"+teamname+"/script.sh")
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(stdout))

		ch <- true
	}()

	return c.SendString(fmt.Sprintf("Received Payload"))

	<-ch

	return c.SendString("Working")

}
