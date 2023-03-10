package controllers

import (
	"fmt"
	"os/exec"

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

	//for non-blocking the script run
	ch := make(chan bool)

	go func() {

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
