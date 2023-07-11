package harbor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	config "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
)

var baseURL string = "https://" + config.KatanaConfig.Harbor.Hostname + "/api/v2.0"

func setAdminPassword() error {
	url := baseURL + "/users/1/password"

	// Make a GET request and read the "X-Harbor-CSRF-Token" header
	// from the response
	resp, err := http.Get(baseURL + "/login")
	if err != nil {
		return err
	}

	csrfToken := resp.Header.Get("X-Harbor-Csrf-Token")

	payload := map[string]interface{}{
		"new_password": config.KatanaConfig.Harbor.Password,
		"old_password": "Harbor12345",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	auth := "Basic " + utils.Base64Encode("admin:Harbor12345") // Harbor12345 is the default admin password

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Harbor-Csrf-Token", csrfToken)
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, _ = client.Do(req)
	if resp.StatusCode != 200 {
		return fmt.Errorf("error changing admin password")
	}

	return nil
}

func createHarborProject(projectName string) error {
	url := baseURL + "/projects"

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
	url := baseURL + "/systeminfo/getcert"

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
