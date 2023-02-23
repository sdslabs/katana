package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// func git(arg string) {
// 	app := "git"
// 	arg0 := arg
// 	cmd := exec.Command(app, arg0)
// 	stdout, err := cmd.Output()

// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	fmt.Println(string(stdout))
// }

//	func Chall() string {
//		cmd, err := exec.Command("/bin/sh", "asdf.sh").Output()
//		if err != nil {
//			fmt.Printf("error %s", err)
//		}
//		output := string(cmd)
//		return output
//	}
func ChallengeUpdate(c *fiber.Ctx) error {
	fmt.Println("Hello")
	//TODO : Apply yaml file and up challenge in container in docker
	return c.SendString("Challenge successfully pulled from gi")
}
