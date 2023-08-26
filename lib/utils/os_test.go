package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestTar(t *testing.T) {
	tmpDir := t.TempDir()
	sampleFileName := "sample_file.txt"
	fileContents := []byte("This is a sample file to test tar function.")
	err := os.WriteFile(tmpDir+string(filepath.Separator)+sampleFileName, fileContents, 0644)
	assert.NoError(t, err)

	buf := new(bytes.Buffer)

	err = Tar(tmpDir, buf)
	assert.NoError(t, err)

	gr, err := gzip.NewReader(buf)
	assert.NoError(t, err)
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		}

		assert.NoError(t, err)

		content, err := io.ReadAll(tr)
		assert.NoError(t, err)

		assert.Equal(t, fileContents, content)

		fileName := strings.TrimPrefix(header.Name, "./"+tmpDir+string(filepath.Separator))
		assert.Equal(t, sampleFileName, fileName)
	}
}

func TestGetKatanaRootPath(t *testing.T) {
	tmpDir := t.TempDir()

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	result, err := GetKatanaRootPath()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if result != tmpDir {
		t.Errorf("Expected result: %s, but got: %s", tmpDir, result)
	}
}
