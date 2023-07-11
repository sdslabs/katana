package harbor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	config "github.com/sdslabs/katana/configs"
	utils "github.com/sdslabs/katana/lib/utils"
)

func addCertificateToDocker() error {
	basePath, _ := os.Getwd()
	cmd := fmt.Sprintf("sudo cp %s/ca.crt /etc/docker/certs.d/"+config.KatanaConfig.Harbor.Hostname+"/ca.crt", basePath)
	out := exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	return nil
}

func setupDockerCredentials() error {
	username := config.KatanaConfig.Harbor.Username
	password := config.KatanaConfig.Harbor.Password
	base64Auth := utils.Base64Encode(username + ":" + password)
	dockerConfig := map[string]map[string]map[string]string{
		"auths": {
			"harbor": {
				"auth": base64Auth,
			},
		},
	}

	dockerConfigJSON, err := json.Marshal(dockerConfig)
	if err != nil {
		return err
	}

	dockerConfigJSONBase64 := utils.Base64Encode(string(dockerConfigJSON))
	dockerConfigJSONBase64 = strings.ReplaceAll(dockerConfigJSONBase64, "\n", "")

	cmd := fmt.Sprintf("echo %s | base64 -d > ~/.docker/config.json", dockerConfigJSONBase64)
	out := exec.Command("bash", "-c", cmd)
	if err := out.Run(); err != nil {
		return err
	}

	return nil
}
