package challengedeployerservice

import (
	"fmt"
	"log"
	"os"
	"regexp"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	v1 "k8s.io/api/core/v1"
)

func CopyInPod(localFilePath string, pathInPod string) error {

	//regex to find challenge name since localFilePath[12:22] is hardcoded
	regexPattern := `\/([^\/]+)\.tar\.gz$`
	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindStringSubmatch(localFilePath)
	filename := matches[1]

	// Get pods from different namespaces
	var pods []v1.Pod
	numberOfTeams := utils.GetTeamNumber()
	for i := 0; i < numberOfTeams; i++ {
		path := "katana-team-" + fmt.Sprint(i) + "/" + filename
		err := os.Mkdir("teams/"+path, 0755)
		if err != nil {
			log.Println(err)
		}
		git.PlainInit("teams/"+path, false)
		repo, err := git.PlainOpen("teams/" + path)
		if err != nil {
			log.Println(err)
		}
		remoteConfig := &config.RemoteConfig{
			Name: "origin",
			URLs: []string{"http://sdslabs@" + utils.GetGogsIp() + ":18080" + "/" + path}}
		_, err = repo.CreateRemote(remoteConfig)

		if err != nil {
			log.Println(err)
		}
		podsInTeam, err := utils.GetPods(map[string]string{
			"app": g.ClusterConfig.TeamLabel,
		}, "katana-team-"+fmt.Sprint(i)+"-ns")
		if err != nil {
			log.Println(err)
			return err
		}
		pods = append(pods, podsInTeam...)
	}
	// Loop over pods
	for _, pod := range pods {
		// Copy file into pod
		if err := utils.CopyIntoPod(pod.Name, g.TeamVmConfig.ContainerName, pathInPod, localFilePath, pod.Namespace); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
