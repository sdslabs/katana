// nolint
package sshproviderservice

// TODO remove nolint later
import (
	"fmt"
	"os"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
)

func CreateTeams() error {
	teamlabels := utils.GetTeamPodLabels()
	var teams []interface{}
	teamnumber := utils.GetTeamNumber()
	credsFile, err := os.Create(g.SSHProviderConfig.CredsFile)
	if err != nil {
		return err
	}
	podName := teamlabels + "-team-master-pod-0"
	for i := 0; i < teamnumber; i++ {
		pwd := utils.GenPassword()
		hashed, err := utils.HashPassword(pwd)
		if err != nil {
			return err
		}
		podNamespace := "katana-team-" + fmt.Sprint(i)
		team := types.CTFTeam{
			Index:     i,
			Name: podNamespace,
			PodName:   podName,
			Password:  hashed,
		}
		fmt.Fprintf(credsFile, "Team: %d, Username: %s, Password: %s\n", i, team.Name, pwd)
		teams = append(teams, team)

		// mysql.CreateGogsUser(team.Namespace, pwd)
		// cmd := exec.Command("kubectl", "exec", "--namespace=", podNamespace, " ", podName, " -- ", "echo", pwd, ">> sshcred")
		// err = cmd.Run()
		// if err != nil {
		// 	panic(err)
		// }

	}

	_, err = mongo.CreateTeams(teams)
	return err
}
