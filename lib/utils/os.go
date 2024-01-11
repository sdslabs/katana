package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(cmd string) error {
	out := exec.Command("bash", "-c", cmd)
	err := out.Run()
	if err != nil {
		return err
	}

	return nil
}

func CreateDirectoryIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(dirPath); err != nil {
			return fmt.Errorf("failed to delete existing directory: %w", err)
		}
		fmt.Printf("Directory '%s' deleted.\n", dirPath)
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	fmt.Printf("Directory '%s' created.\n", dirPath)
	return nil
}
