package infraService

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	g "github.com/sdslabs/katana/configs"
)

// runCmd represents the run command
// not tested yet [WIP]
var SetUpCmd = &cobra.Command{
	Use:   "setup",
	Short: "SetUps Katana on your computer",
	RunE: func(cmd *cobra.Command, args []string) error {
		g.LoadConfiguration();
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
		time.Sleep(5 * time.Second)
		if err := GitSetup(); err != nil {
			log.Println("Error setting up the git server:", err)
			return err
		}
		log.Println("Git Server connected successfully")
		time.Sleep(5 * time.Second)
		if err := mongoDBSetup(); err != nil {
			log.Println("Error setting up the mysql database:", err)
			return err
		}
		time.Sleep(5 * time.Second)
		if err := mysqlDBSetup(); err != nil {
			log.Println("Error setting up the mysql database:", err)
			return err
		}
		log.Println("MySQL Database connected successfully")
		return nil
	},
}
