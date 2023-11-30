package chalDeployerService

import (
	"log"

	"github.com/spf13/cobra"
)

var DeployChalCmd = &cobra.Command{
	Use:   "chal-deploy",
	Short: "Run the Challenge Update command",
	Long:  "Runs the challenge update",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := DeployChallenge(); err != nil {
			log.Println("Error deploying the challenge:", err)
			return err
		}
		log.Println("Challenge deployed successfully")
		return nil
	},
}
