package harbor

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	config "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
)

var baseURL string = "https://harbor.katana.local/api/v2.0"

var httpClient *http.Client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func setAdminPassword() error {
	url := baseURL + "/users/1/password"

	// Make a GET request and read the "X-Harbor-CSRF-Token" header
	// from the response
	resp, err := httpClient.Get(baseURL)
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

	resp, _ = httpClient.Do(req)
	if resp.StatusCode != 200 {
		return fmt.Errorf("error changing admin password, response code is %d", resp.StatusCode)
	}

	return nil
}

func createHarborProject(projectName string) error {
	url := baseURL + "/projects"

	payload := []byte(`{
		"project_name": "katana",
		"metadata": {
			"public": "true"
		},
		"storage_limit": -1
	}`)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	auth := "Basic " + utils.Base64Encode("admin:"+config.KatanaConfig.Harbor.Password)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	resp, _ := httpClient.Do(req)
	if resp.StatusCode != 201 {
		return fmt.Errorf("error creating project")
	}

	return nil
}
