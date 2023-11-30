package chalDeployerService

import (
	"log"

	"github.com/spf13/cobra"
)

var DelChalCmd = &cobra.Command{
	Use:   "delete-chal",
	Short: "Run the Challenge Delete command",
	Long:  "Deletes the Challenge",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := DeleteChallenge(args[0]); err != nil {
			log.Println("Error deleting the challenge:", err)
			return err
		}
		log.Println("Challenge deleted successfully")
		return nil
	},
}
