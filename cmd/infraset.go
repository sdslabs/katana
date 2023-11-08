package main

import (
	"log"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/deployment"
	"github.com/sdslabs/katana/lib/harbor"
	utils "github.com/sdslabs/katana/lib/utils"
	infraSetService "github.com/sdslabs/katana/services/infrasetservice"
	"github.com/spf13/cobra"
)

var infraCmd = &cobra.Command{

	Use:   "infra-setup",
	Short: "Run the infraset setup command",
	Long:  `Runs the katana API server on port 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		g.LoadConfiguration()
		config, err := utils.GetKubeConfig()

		if err != nil {
			log.Println("There was error in setting up Kubernentes Config", err)
		}

		kubeclient, err := utils.GetKubeClient()
		if err != nil {
			log.Println("Error in creating Kubernetes client", err)
		}
		infraSetService.GenerateCertsforHarbor()
		deployment.DeployCluster(config, kubeclient)
		err = harbor.SetupHarbor()
		if err != nil {
			log.Println("There was error in setting up harbor", err)
		}
		infraSetService.BuildKatanaServices()
	},
}
