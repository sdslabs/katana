package controllers

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/gofiber/fiber/v2"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/utils"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

/* func CreateTeam(namespace string, i string) string {
	// ctfTeam := new(types.CTFTeam)
	log.Println("Creating pod for team " + i)

	return "PlaceHolder"
} */

func CreateTeams(c *fiber.Ctx) error {

	// if !utils.VerifyToken(c) {
	// 	return c.SendString("Unauthorized")
	// }

	clusterConfig := g.ClusterConfig
	client, err := utils.GetKubeClient()
	if err != nil {
		log.Println(err)
	}
	noOfTeams, err := strconv.Atoi(c.Params("number"))

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < noOfTeams; i++ {
		log.Println("Creating Team: " + strconv.Itoa(i))
		nsName := &coreV1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "katana-team-ns-" + strconv.Itoa(i),
			},
		}
		_, err = client.CoreV1().Namespaces().Create(c.Context(), nsName, metav1.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}
		manifest := &bytes.Buffer{}
		tmpl, err := template.ParseFiles(filepath.Join(clusterConfig.ManifestDir, "teams.yml"))
		if err != nil {
			return err
		}
		deploymentConfig := utils.DeploymentConfig()

		if err = tmpl.Execute(manifest, deploymentConfig); err != nil {
			return err
		}
		pathToCfg := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
		config, err := clientcmd.BuildConfigFromFlags("", pathToCfg)
		if err != nil {
			log.Fatal(err)
		}
		if err = deployment.ApplyManifest(config, client, manifest.Bytes(), "katana-team-ns-"+strconv.Itoa(i)); err != nil {
			return err
		}
	}
	return c.SendString("Successfully created teams")
}
