package main

import (
	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/services/master"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the katana API server",
	Long:  `Runs the katana API server on port 3000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// var g errgroup.Group

		// vmDeployerListener, err := net.Listen("tcp", ":8001")
		// if err != nil {
		// 	panic(err)
		// }

		// vmDeployerServer := vmdeployerservice.Server()

		// if e := vmDeployerServer.Serve(vmDeployerListener); e != nil {
		// 	panic(e)
		// }

		// apiServer := api.Server()
		configs.LoadConfiguration()
		return master.Server()

		// if err := g.Wait(); err != nil {
		// 	os.Exit(1)
		// }
	},
}
