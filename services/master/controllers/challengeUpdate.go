package controllers

import (
	"fmt"
	"os/exec"
	"github.com/gofiber/fiber/v2"
)


type Repo struct {
	Name string `json:"name"`
}

type Teams struct {
	Ref  string `json:"ref"`
	Before   string `json:"before"`
	Repository Repo `json:"repository"`
}

func ChallengeUpdate(c *fiber.Ctx) error {

	
	p := new(Teams)
	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse JSON request body: %v", err))
	}

	// fmt.Println(p.Repository)
	// fmt.Println(p.Repository.Name)
	fmt.println("Request Recieved, please wait..pulling repo and building images")
	teamname := p.Repository.Name
	//fmt.Println(teamname)
	
	//for non-blocking the script run
	ch := make(chan bool)

	go func(){

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
