package main

import (
	"log"

	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{

	Use:   "git-server",
	Short: "Run the git-server setup command",
	Long:  `Runs the katana API server on port 3000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gitSetup(); err != nil {
			log.Println("Error setting up the git server:", err)
			return err
		}
		log.Println("Git Server connected successfully")
		return nil
	},
}
