package main

import (
	"log"

	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db-setup",
	Short: "Run the db setup command",
	Long:  "Runs the database setup",
	RunE: func(cmd *cobra.Command, args []string) error{
		if err := dbSetup(); err != nil {
			log.Println("Error setting up the database:", err)
			return err
		}
		log.Println("Database connected successfully")
		return nil
	},
}
