package main

import (
	"log"

	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/mysql"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{

	Use:   "db-setup",
	Short: "Run the db setup command",
	Long:  `Runs the katana API server on port 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := mongo.Init(); err != nil {
			log.Println("There was error in settong up mongo db", err)
		}
		if err := mysql.Init(); err != nil {
			log.Println("There was error in settong up mysql db", err)
		}
		log.Println("DB Connected Succesfully")
	},
}
