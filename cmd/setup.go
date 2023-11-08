package main

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var setUpCmd = &cobra.Command{
	Use:   "setup",
	Short: "SetUps Katana on your computer",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
