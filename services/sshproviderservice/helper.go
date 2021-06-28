package sshproviderservice

import (
	"fmt"
	"os"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/utils"
)

func createTeams() error {
	teamlabels := make(map[string]string)
	teamlabels["app"] = g.ClusterConfig.TeamLabel
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

		team := CTFTeam{
			Index:    i,
			Name:     pod.Name,
			PodName:  pod.Name,
			Password: hashed,
		}
		fmt.Fprintf(credsFile, "Team: %d, Username: %s, Password: %s\n", i, team.Name, pwd)
		teams = append(teams, team)
	}

	_, err = mongo.CreateTeams(teams)
	return err
}
