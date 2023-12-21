package cliHelpers

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"

	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/services/infrasetservice"
)

func CreateTeams() error {

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
		_, err = client.CoreV1().Namespaces().Create(context.TODO(), nsName, metav1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
			return err
		}
		
		manifest := &bytes.Buffer{}
		tmpl, err := template.ParseFiles(filepath.Join(configs.ClusterConfig.TemplatedManifestDir, "runtime", "teams.yml"))
		if err != nil {
			return err
		}

		pwd, team, err := infrasetservice.CreateTeamCredentials(i)
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
	return nil
}
