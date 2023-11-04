package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path/filepath"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/harbor"
	utils "github.com/sdslabs/katana/lib/utils"
	infraSetService "github.com/sdslabs/katana/services/infrasetservice"
	"github.com/spf13/cobra"
)

var infraCmd = &cobra.Command{

	Use:   "infra",
	Short: "Run the infraset setup command",
	Long:  `Runs the katana API server on port 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		g.LoadConfiguration()
		config, err := utils.GetKubeConfig()

		if err != nil {
			log.Println("There was error in setting up Kubernentes Config")
			// return err
		}

		kubeclient, err := utils.GetKubeClient()

		if err != nil {
			log.Fatal(err)
		}
		infraSetService.GenerateCertsforHarbor()

		clusterConfig := g.ClusterConfig
		deploymentConfig := utils.DeploymentConfig()
		nodes, _ := utils.GetNodes(kubeclient)

		deploymentConfig.NodeAffinityValue = nodes[0].Name

		for _, m := range clusterConfig.TemplatedManifests {

			manifest := &bytes.Buffer{}
			log.Printf("Applying: %s\n", m)
			tmpl, err := template.ParseFiles(filepath.Join(clusterConfig.TemplatedManifestDir, m))
			if err != nil {
				fmt.Print("err")
			}

			if err = tmpl.Execute(manifest, deploymentConfig); err != nil {
				fmt.Print("err")
			}

			if err = deployment.ApplyManifest(config, kubeclient, manifest.Bytes(), g.KatanaConfig.KubeNameSpace); err != nil {
				fmt.Print("err")
			}
		}
		err = harbor.SetupHarbor()
		if err != nil {
			fmt.Print("err")
		}
		infraSetService.BuildKatanaServices()
	},
}
