package main

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the run command
var rootCmd = &cobra.Command{
	Use:   "katana",
	Short: "Katana is an automatic cloud-native attack-defence CTF platform",
	Long: `Katana is an automatic cloud-native attack-defence CTF platform 
	that streamlines the setup and management of attack-defence CTFs with automated infrastructure for
	challenge dispatcher, VM deployer and flag juggler.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
