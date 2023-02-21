package main

import (
	"fmt"
	"os/exec"
)

func git(arg string) {
	app := "git"
	arg0 := arg
	cmd := exec.Command(app, arg0)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

func main() {
	git("pull")
}
