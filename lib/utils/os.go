package utils

import "os/exec"

func RunCommand(cmd string) error {
	out := exec.Command("bash", "-c", cmd)
	err := out.Run()
	if err != nil {
		return err
	}

	return nil
}
