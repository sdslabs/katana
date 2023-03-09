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

func CreateTeams() error {
	teamlabels := utils.GetTeamPodLabels()
	var teams []interface{}
	teamnumber := utils.GetTeamNumber()
	credsFile, err := os.Open(g.SSHProviderConfig.CredsFile)
	if err != nil {
		return err
	}
	for i := 0; i < teamnumber; i++ {
		pwd := utils.GenPassword()
		hashed, err := utils.HashPassword(pwd)
		if err != nil {
			return err
		}
		podName := teamlabels + "-teams" + fmt.Sprint(teamnumber)
		team := types.CTFTeam{
			Index:    i,
			Name:     podName,
			PodName:  podName,
			Password: hashed,
		}
		fmt.Fprintf(credsFile, "Team: %d, Username: %s, Password: %s\n", i, team.Name, pwd)
		teams = append(teams, team)

		mysql.CreateGogsUser(team.Name, pwd)
		cmd := exec.Command("kubectl exec ", podName, " -- touch sshcreds")
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		cmd = exec.Command("kubectl exec ", podName, " -- sh -c 'echo", pwd, "' >> sshcred")
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

	}

	_, err = mongo.CreateTeams(teams)
	return err
}
