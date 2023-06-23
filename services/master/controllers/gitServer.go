package controllers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GitServer(c *fiber.Ctx) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	kubeclient, err := utils.GetKubeClient("")
	if err != nil {
		log.Println(err)
	}

	service, err := kubeclient.CoreV1().Services("katana").Get(context.TODO(), "mysql-nodeport-svc", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}

	mysqlAddress := service.Spec.ClusterIP
	gogsService, err := kubeclient.CoreV1().Services("katana").Get(context.TODO(), "gogs-svc", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}

	gitServerAddress := gogsService.Spec.ClusterIP + ":18080"

	// Add the form fields
	writer.WriteField("db_type", "MySQL")
	writer.WriteField("db_host", mysqlAddress)
	writer.WriteField("db_user", configs.MySQLConfig.Username)
	writer.WriteField("db_passwd", configs.MySQLConfig.Password)
	writer.WriteField("db_name", "gogs")
	writer.WriteField("db_schema", "public")
	writer.WriteField("ssl_mode", "disable")
	writer.WriteField("db_path", "/app/gogs/data/gogs.db")
	writer.WriteField("app_name", "Gogs")
	writer.WriteField("repo_root_path", "/data/git/gogs-repositories")
	writer.WriteField("run_user", "git")
	writer.WriteField("domain", gitServerAddress)
	writer.WriteField("ssh_port", "22")
	writer.WriteField("http_port", "3000")
	writer.WriteField("app_url", "http://"+gitServerAddress+":3000/")
	writer.WriteField("log_root_path", "/app/gogs/log")
	writer.WriteField("default_branch", "master")

	// Close the writer
	writer.Close()

	// Create the request
	req, err := http.NewRequest("POST", "http://"+gitServerAddress+"/install", &requestBody)
	if err != nil {
		return err
	}

	// Set the content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error while setting up Git Server")
	}
	
	return c.SendString("Git Server setup completed\n")
}
