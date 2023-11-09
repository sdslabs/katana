package main

import (
	"log"

	"github.com/spf13/cobra"
)

var infraCmd = &cobra.Command{

	Use:   "infra-setup",
	Short: "Run the infraset setup command",
	Long:  `Runs the katana API server on port 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := infraSetup(); err != nil {
			log.Println("Error setting up the infrastructure:", err)
		}
		log.Println("Infrastructure setup successfully")
	},
}
