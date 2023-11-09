package main

import (
	"log"

	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{

	Use:   "git-server",
	Short: "Run the git-server setup command",
	Long:  `Runs the katana API server on port 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gitSetup(); err != nil {
			log.Println("Error setting up the git server:", err)
			return 
		}
		log.Println("Git Server connected successfully")
	},
}
