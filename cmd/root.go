package main

import (
	"log"

	"github.com/sdslabs/katana/cliHelpers/chalDeployerService"
	"github.com/sdslabs/katana/cliHelpers/infraService"
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
		if err := cmd.Help(); err != nil {
			log.Printf("Failed to print cobra help: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(infraService.MongoDBCmd)
	rootCmd.AddCommand(infraService.MysqlDBCmd)
	rootCmd.AddCommand(infraService.InfraCmd)
	rootCmd.AddCommand(infraService.SetUpCmd)
	rootCmd.AddCommand(infraService.GitCmd)
	rootCmd.AddCommand(chalDeployerService.DelChalCmd)
	rootCmd.AddCommand(chalDeployerService.ChalUpdateCmd)
	rootCmd.AddCommand(chalDeployerService.DeployChalCmd)
}
