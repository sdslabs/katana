package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/sdslabs/katana/configs"
	utils "github.com/sdslabs/katana/lib/utils"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{

	Use:   "git-server",
	Short: "Run the git-server setup command",
	Long:  `Runs the katana API server on port 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)
		LoadBalancer := utils.GetKatanaLoadbalancer()

		writer.WriteField("db_type", "MySQL")
		writer.WriteField("db_host", LoadBalancer+":3306")
		writer.WriteField("db_user", configs.MySQLConfig.Username)
		writer.WriteField("db_passwd", configs.MySQLConfig.Password)
		writer.WriteField("db_name", "gogs")
		writer.WriteField("db_schema", "public")
		writer.WriteField("ssl_mode", "disable")
		writer.WriteField("db_path", "/app/gogs/data/gogs.db")
		writer.WriteField("app_name", "Gogs")
		writer.WriteField("repo_root_path", "/data/git/gogs-repositories")
		writer.WriteField("run_user", "git")
		writer.WriteField("domain", LoadBalancer+":3000")
		writer.WriteField("ssh_port", "22")
		writer.WriteField("http_port", "3000")
		writer.WriteField("app_url", "http://"+LoadBalancer+":3000")
		writer.WriteField("log_root_path", "/app/gogs/log")
		writer.WriteField("default_branch", "master")

		// Close the writer
		writer.Close()

		// Create the request
		req, err := http.NewRequest("POST", "http://"+LoadBalancer+":3000"+"/install", &requestBody)
		if err != nil {
			fmt.Println("", err)
		}

		// Set the content type
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("",err)
		}

		// Check the response
		if resp.StatusCode != http.StatusOK {
			fmt.Println("error while setting up Git Server")
		}

		log.Println("Git Server setup completed")

	},
}
