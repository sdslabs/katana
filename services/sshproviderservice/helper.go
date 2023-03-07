// nolint
package sshproviderservice

// TODO remove nolint later
import (
	"fmt"
	"os"
	"os/exec"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/mysql"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
)

func createTeams() error {
	teamlabels := utils.GetTeamPodLabels()
	var teams []interface{}
	teamPods, err := utils.GetPods(kubeClientset, teamlabels)
	if err != nil {
		return err
	}

	credsFile, err := os.Open(g.SSHProviderConfig.CredsFile)
	if err != nil {
		return err
	}

	for i, pod := range teamPods {
		pwd := utils.GenPassword()
		hashed, err := utils.HashPassword(pwd)
		if err != nil {
			return err
		}

		team := types.CTFTeam{
			Index:    i,
			Name:     pod.Name,
			PodName:  pod.Name,
			Password: hashed,
		}
		fmt.Fprintf(credsFile, "Team: %d, Username: %s, Password: %s\n", i, team.Name, pwd)
		teams = append(teams, team)

		mysql.CreateGogsUser(team.Name, pwd)
		cmd := exec.Command("kubectl exec ",pod.Name," -- touch sshcreds")
		err = cmd.Run()
		if(err != nil){
			panic(err)
		}
		cmd = exec.Command("kubectl exec ",pod.Name," -- sh -c 'echo", pwd,"' >> sshcred")
		err = cmd.Run()
		if(err != nil){
			panic(err)
		}

	}

	_, err = mongo.CreateTeams(teams)
	return err
}
