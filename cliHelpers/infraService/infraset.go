package infraService

import (
	"log"

	"github.com/spf13/cobra"

	g "github.com/sdslabs/katana/configs"
)

var InfraCmd = &cobra.Command{

	Use:   "infra-setup",
	Short: "Run the infraset setup command",
	Long:  `Runs the katana API server on port 3000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		g.LoadConfiguration();
		if err := InfraSetup(); err != nil {
			log.Println("Error setting up the infrastructure:", err)
			return err
		}
		log.Println("Infrastructure setup successfully")
		return nil
	},
}
