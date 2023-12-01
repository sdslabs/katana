package chalDeployerService

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var DelChalCmd = &cobra.Command{
	Use:   "delete-chal",
	Short: "Run the Challenge Delete command",
	Long:  "Deletes the Challenge",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		challengeID := args[0]

		if err := DeleteChallenge(challengeID); err != nil {
			log.Printf("Error deleting the challenge with ID %s: %v", challengeID, err)
			return fmt.Errorf("failed to delete challenge: %v", err)
		}

		log.Printf("Challenge with ID %s deleted successfully", challengeID)
		return nil
	},
}
