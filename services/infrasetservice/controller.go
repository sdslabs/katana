package infrasetservice

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/harbor"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/mysql"
	utils "github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
)

func InfraSet(c *fiber.Ctx) error {
	config:= configs.GlobalKubeConfig
	kubeclient:= configs.GlobalKubeClient
	if err := GenerateCertsforHarbor(); err != nil {
		log.Println("Error generating certificates :", err)
		return err
	}
	if err := deployment.DeployCluster(config, kubeclient); err != nil {
		log.Fatal(err)
		return err
	}
	err := harbor.SetupHarbor()
	if err != nil {
		log.Fatal(err)
		return err
	}
	if err = BuildKatanaServices(); err != nil {
		log.Fatal(err)
		return err
	}

	return c.SendString("Infrastructure setup completed")
}

func DB(c *fiber.Ctx) error {
	// TODO: run Mongo and MySQL setup in parallel
	if err := mongo.Init(); err != nil {
		return err
	}
	if err := mysql.Init(); err != nil {
		return err
	}
	return c.SendString("Database setup completed\n")
}

func Login(c *fiber.Ctx) error { // TODO: Remove this function from this package and move it to its own package
	adminUser := new(types.AdminUser)

	if err := c.BodyParser(adminUser); err != nil {
		return err
	}

	admin, err := mongo.FetchSingleAdmin(adminUser.Username)
	if err != nil {
		return err
	}

	isAdmin := utils.CompareHashWithPassword(admin.Password, adminUser.Password)
	if isAdmin {
		jwtToken := jwt.New(jwt.SigningMethodHS256)

		token, err := jwtToken.SignedString([]byte("secret"))
		if err != nil {
			return err
		}

		cookie := new(fiber.Cookie)
		cookie.Name = "jwt"
		cookie.Value = token
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.HTTPOnly = true
		c.Cookie(cookie)

		return c.JSON(fiber.Map{
			"message": "success",
		})
	} else {
		return c.SendString("Incorrect password")
	}
}

func CreateTeams(c *fiber.Ctx) error {

	// if !utils.VerifyToken(c) {
	// 	return c.SendString("Unauthorized")
	// }
	config:= configs.GlobalKubeConfig
	client:= configs.GlobalKubeClient

	noOfTeams := int(configs.ClusterConfig.TeamCount)

	if _, err := os.Stat("teams"); os.IsNotExist(err) {
		errDir := os.Mkdir("teams", 0755)
		if errDir != nil {
			log.Fatal(err)
		return err
		}
	}

	// Create a directory named teams in the current directory
	if _, err := os.Stat("teams"); os.IsNotExist(err) {
		errDir := os.Mkdir("teams", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	// Create a directory named teams in the current directory
	if _, err := os.Stat("teams"); os.IsNotExist(err) {
		errDir := os.Mkdir("teams", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	var teams []interface{}
	credsFile, err := os.Create(configs.SSHProviderConfig.CredsFile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for i := 0; i < noOfTeams; i++ {
		// Create a directory named katana-team-i in the teams directory
		if _, err := os.Stat("teams/katana-team-" + strconv.Itoa(i)); os.IsNotExist(err) {
			errDir := os.Mkdir("teams/katana-team-"+strconv.Itoa(i), 0755)
			if errDir != nil {
				log.Fatal(err)
				return err
			}
		}

		log.Println("Creating Team: " + strconv.Itoa(i))
		namespace := "katana-team-" + strconv.Itoa(i) + "-ns"
		nsName := &coreV1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		// var contx context.Context
		_, err = client.CoreV1().Namespaces().Create(c.Context(), nsName, metav1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
			return err
		}
		manifest := &bytes.Buffer{}
		tmpl, err := template.ParseFiles(filepath.Join(configs.ClusterConfig.TemplatedManifestDir, "runtime", "teams.yml"))
		if err != nil {
			log.Fatal(err)
			return err
		}
		pwd, team, err := CreateTeamCredentials(i)
		if err != nil {
			return err
		}
		deploymentConfig := utils.DeploymentConfig()

		deploymentConfig.SSHPassword = pwd
		if err = tmpl.Execute(manifest, deploymentConfig); err != nil {
			return err
		}

		if err = deployment.ApplyManifest(config, client, manifest.Bytes(), namespace); err != nil {
			return err
		}
		teams = append(teams, team)
		fmt.Fprintf(credsFile, "Team: %d, Username: %s, Password: %s\n", i, team.Name, pwd)
	}
	mongo.CreateTeams(teams)
	return c.SendString("Successfully created teams")
}

func GitServer(c *fiber.Ctx) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	LoadBalancer, err := utils.GetKatanaLoadbalancer()
	if err != nil {
		return err
	}
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
