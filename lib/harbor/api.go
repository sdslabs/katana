package harbor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	config "github.com/sdslabs/katana/configs"
)

func createHarborProject(projectName string) error {
	url := "https://" + config.KatanaConfig.Harbor.Hostname + "/api/v2.0/projects"

	payload := map[string]interface{}{
		"name": projectName,
		"metadata": map[string]interface{}{
			"public": "true",
		},
		"storage_limit": -1,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	if resp.StatusCode != 201 {
		return fmt.Errorf("error creating project")
	}

	return nil
}

func getHarborCertificate() error {
	url := "https://" + config.KatanaConfig.Harbor.Hostname + "/api/v2.0/systeminfo/getcert"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/octet-stream")

	client := &http.Client{}
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		return fmt.Errorf("error getting certificate")
	}

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)

	// Save certificate to file
	basePath, _ := os.Getwd()
	certFile, err := os.Create(fmt.Sprintf("%s/ca.crt", basePath))
	if err != nil {
		return err
	}

	defer certFile.Close()

	if _, err := certFile.Write([]byte(respBody["certificate"].(string))); err != nil {
		return err
	}

	return nil
}
