package infraService

import (
	"log"

	"github.com/spf13/cobra"
)

var MongoDBCmd = &cobra.Command{
	Use:   "mongodb-setup",
	Short: "Run the mongo db setup command",
	Long:  "Runs the mongo database setup",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := mongoDBSetup(); err != nil {
			log.Println("Error setting up the database:", err)
			return err
		}
		log.Println("Database connected successfully")
		return nil
	},
}
var MysqlDBCmd = &cobra.Command{
	Use:   "mysqldb-setup",
	Short: "Run the mysql db setup command",
	Long:  "Runs the mysql database setup",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := mysqlDBSetup(); err != nil {
			log.Println("Error setting up the database:", err)
			return err
		}
		log.Println("Database connected successfully")
		return nil
	},
}
