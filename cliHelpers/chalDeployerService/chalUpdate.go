package chalDeployerService

import (
	"log"

	"github.com/spf13/cobra"
)

var ChalUpdateCmd = &cobra.Command{
	Use:   "chal-update",
	Short: "Run the Challenge Update command",
	Long:  "Runs the challenge update",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ChallengeUpdate(); err != nil {
			log.Println("Error updating the challenge:", err)
			return err
		}
		log.Println("Challenge updated successfully")
		return nil
	},
}
