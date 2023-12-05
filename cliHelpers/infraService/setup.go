package infraService

import (
	"log"

	"github.com/spf13/cobra"
)

// runCmd represents the run command 
// not tested yet [WIP]
var SetUpCmd = &cobra.Command{
	Use:   "setup",
	Short: "SetUps Katana on your computer",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := InfraSetup(); err != nil {
			log.Println("Error setting up the infrastructure:", err)
			return err
		}
		log.Println("Infrastructure setup successfully")
		if err := mongoDBSetup(); err != nil {
			log.Println("Error setting up the mongo database:", err)
			return err
		}
		log.Println("MongoDB connected successfully")
		if err := GitSetup(); err != nil {
			log.Println("Error setting up the git server:", err)
			return err
		}
		log.Println("Git Server connected successfully")
		if err := mysqlDBSetup(); err != nil {
			log.Println("Error setting up the mysql database:", err)
			return err
		}
		log.Println("MySQL Database connected successfully")
		return nil
	},
}
