package infraService

import (
"fmt"
"log"

	"github.com/spf13/cobra"

	g "github.com/sdslabs/katana/configs"
"github.com/sdslabs/katana/cliHelpers"
)

var GitCmd = &cobra.Command{
	Use:   "git-server",
	Short: "Run the git-server setup command",
	Long:  `Runs the katana API server on port 3000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		g.LoadConfiguration();
		g.LoadKubeElements();
		if err := GitSetup(); err != nil {
			log.Println("Error setting up the git server:", err)
			return err
		}
		log.Println("Git Server connected successfully")
		return nil
	},
}

var CreateTeamCmd= &cobra.Command{	
	Use:   "create-team",
	Short: "Run the create-team command",
	Long:  `Create namespaces fot teams`,
	RunE: func(cmd *cobra.Command, args []string) error {
		g.LoadConfiguration();
		g.LoadKubeElements();
		fmt.Println("Creating teams")
		if err := cliHelpers.CreateTeams(); err != nil {
			log.Println("Error creating teams:", err)
			return err
		}
		return nil
	},
}