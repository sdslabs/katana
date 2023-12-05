package chalDeployerService

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	g "github.com/sdslabs/katana/configs"
)
// still have to test this [WIP]

var DelChalCmd = &cobra.Command{
	Use:   "delete-chal",
	Short: "Run the Challenge Delete command",
	Long:  "Deletes the Challenge",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		challengeID := args[0]
		g.LoadConfiguration();
		g.LoadKubeElements();

		if err := DeleteChallenge(challengeID); err != nil {
			log.Printf("Error deleting the challenge with ID %s: %v", challengeID, err)
			return fmt.Errorf("failed to delete challenge: %v", err)
		}

		log.Printf("Challenge with ID %s deleted successfully", challengeID)
		return nil
	},
}
