package chalDeployerService

import (
	"log"

	"github.com/spf13/cobra"
)

var DeployChalCmd = &cobra.Command{
	Use:   "chal-deploy",
	Short: "Run the Challenge Deploy command",
	Long:  "Runs the challenge deploy",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := DeployChallenge(); err != nil {
			log.Println("Error deploying the challenge:", err)
			return err
		}
		log.Println("Challenge deployed successfully")
		return nil
	},
}
