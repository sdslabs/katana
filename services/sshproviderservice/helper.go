// nolint
package sshproviderservice

// TODO remove nolint later
import (
	"fmt"
	"os"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/mysql"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
)

func CreateTeams(teamnumber int) error {
	teamlabels := utils.GetTeamPodLabels()
	var teams []interface{}
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
			Index:    i,
			Name:     podNamespace,
			PodName:  podName,
			Password: hashed,
		}
		fmt.Fprintf(credsFile, "Team: %d, Username: %s, Password: %s, Hash: %s\n", i, team.Name, pwd, hashed)
		teams = append(teams, team)
		mysql.CreateGogsUser(team.Name, pwd)

	}
	_, err = mongo.CreateTeams(teams)
	return err
}
