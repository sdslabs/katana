package infraService

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/harbor"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/mysql"
	utils "github.com/sdslabs/katana/lib/utils"
	infraSetService "github.com/sdslabs/katana/services/infrasetservice"
)

func mongoDBSetup() error {
	if err := mongo.Init(); err != nil {
		return err
	}
	return nil
}
func mysqlDBSetup() error {
	if err := mysql.Init(); err != nil {
		return err
	}
	return nil
}

func InfraSetup() error {
	config:= g.GlobalKubeConfig
	kubeclient:= g.GlobalKubeClient

	infraSetService.GenerateCertsforHarbor()
	deployment.DeployCluster(config, kubeclient)
	err := harbor.SetupHarbor()
	if err != nil {
		log.Println("There was error in setting up harbor", err)
		return err
	}
	err = infraSetService.BuildKatanaServices()
	if err != nil {
		return err
	}
	return nil
}

func GitSetup() error {
	if err := mysql.Init(); err != nil {
		return err
	}
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	LoadBalancer, err := utils.GetKatanaLoadbalancer()
	if err != nil {
		return err
	}
	writer.WriteField("db_type", "MySQL")
	writer.WriteField("db_host", LoadBalancer+":3306")
	writer.WriteField("db_user", g.MySQLConfig.Username)
	writer.WriteField("db_passwd", g.MySQLConfig.Password)
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
		return err
	}

	// Set the content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("", err)
		return err
	}

	// Check the response
	if resp.StatusCode != http.StatusOK {
		fmt.Println("error while setting up Git Server")
		return err
	}

	log.Println("Git Server setup completed")
	return nil
}
