package harbor

import (
	"fmt"
	"os"
	"os/exec"

	config "github.com/sdslabs/katana/configs"
)

func setCertificateToDocker() error {
	path := "/etc/docker/certs.d/"

	// Make the directory if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return err
		}
	}

	path = path + config.KatanaConfig.Harbor.Hostname

	// Make the directory if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return err
		}
	}

	basePath, _ := os.Getwd()
	cmd := fmt.Sprintf("sudo cp %s/lib/harbor/certs/ca.crt /etc/docker/certs.d/"+config.KatanaConfig.Harbor.Hostname+"/ca.crt", basePath)
	out := exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	return nil
}

func dockerLogin() error {
	username := config.KatanaConfig.Harbor.Username
	password := config.KatanaConfig.Harbor.Password

	cmd := fmt.Sprintf("sudo docker login -u %s -p %s %s", username, password, config.KatanaConfig.Harbor.Hostname)
	out := exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	return nil
}
