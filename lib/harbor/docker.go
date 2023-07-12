package harbor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sdslabs/katana/configs"
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

	path = path + configs.KatanaConfig.Harbor.Hostname

	// Make the directory if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return err
		}
	}

	basePath, _ := os.Getwd()
	cmd := fmt.Sprintf("sudo cp %s/lib/harbor/certs/ca.crt /etc/docker/certs.d/"+configs.KatanaConfig.Harbor.Hostname+"/ca.crt", basePath)
	out := exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	return nil
}
