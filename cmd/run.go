package main

import (
	"github.com/sdslabs/katana/api"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the katana API server",
	Long:  `Runs the katana API server on port 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		api.RunKatanaAPIServer()
	},
}
