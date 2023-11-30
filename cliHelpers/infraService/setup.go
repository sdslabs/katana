package infraService

import (
	"log"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var SetUpCmd = &cobra.Command{
	Use:   "setup",
	Short: "SetUps Katana on your computer",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := InfraSetup(); err != nil {
			log.Println("Error setting up the infrastructure:", err)
			return err
		}
		log.Println("Infrastructure setup successfully")
		if err := DBSetup(); err != nil {
			log.Println("Error setting up the database:", err)
			return err
		}
		log.Println("Database connected successfully")
		if err := GitSetup(); err != nil {
			log.Println("Error setting up the git server:", err)
			return err
		}
		log.Println("Git Server connected successfully")
		if err := DBSetup(); err != nil {
			log.Println("Error setting up the database:", err)
			return err
		}
		log.Println("Database connected successfully")
		return nil
	},
}
