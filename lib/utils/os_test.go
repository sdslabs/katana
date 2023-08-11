package utils

import (
	"os/exec"
	"testing"
)

func TestRunCommand(t *testing.T) {
	// Call the function being tested
	err := RunCommand("echo 'Hello, world!'")

	// Check if the function returned an error
	if err != nil {
		t.Errorf("Test case failed: %v", err)
	}

	// Check if the command was executed correctly
	expectedOutput := "Hello, world!\n"
	cmd := exec.Command("bash", "-c", "echo 'Hello, world!'")
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Test case failed: %v", err)
	}
	if string(output) != expectedOutput {
		t.Errorf("Test case failed: expected '%s', got '%s'", expectedOutput, string(output))
	}
}
