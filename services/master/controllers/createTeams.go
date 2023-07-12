package controllers

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/utils"
	ssh "github.com/sdslabs/katana/services/sshproviderservice"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateTeams(c *fiber.Ctx) error {

	// if !utils.VerifyToken(c) {
	// 	return c.SendString("Unauthorized")
	// }

	config, err := utils.GetKubeConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err := utils.GetKubeClient()
	if err != nil {
		log.Fatal(err)
	}
	noOfTeams, err := strconv.Atoi(c.Params("number"))

	if err != nil {
		log.Fatal(err)
	}
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

	// Create a directory named teams in the current directory
	if _, err := os.Stat("teams"); os.IsNotExist(err) {
		errDir := os.Mkdir("teams", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	for i := 0; i < noOfTeams; i++ {
		// Create a directory named katana-team-i in the teams directory
		if _, err := os.Stat("teams/katana-team-" + strconv.Itoa(i)); os.IsNotExist(err) {
			errDir := os.Mkdir("teams/katana-team-"+strconv.Itoa(i), 0755)
			if errDir != nil {
				log.Fatal(err)
			}
		}

		log.Println("Creating Team: " + strconv.Itoa(i))
		namespace := "katana-team-" + strconv.Itoa(i) + "-ns"
		nsName := &coreV1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		_, err = client.CoreV1().Namespaces().Create(c.Context(), nsName, metav1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}
		manifest := &bytes.Buffer{}
		tmpl, err := template.ParseFiles(filepath.Join(g.ClusterConfig.ManifestRuntimeDir, "teams.yml"))
		if err != nil {
			return err
		}
		deploymentConfig := utils.DeploymentConfig()

		if err = tmpl.Execute(manifest, deploymentConfig); err != nil {
			return err
		}

		if err = deployment.ApplyManifest(config, client, manifest.Bytes(), namespace); err != nil {
			return err
		}
	}
	SSH(noOfTeams)
	return c.SendString("Successfully created teams")
}

func SSH(noOfTeams int) {
	ssh.CreateTeams(noOfTeams)
	startServer()
}

func startServer() {
	x := ssh.Server()
	go func() {
		x.ListenAndServe()
	}()
	log.Println("Server up and running")
	for {
		//to keeep this server running forever
	}
}
