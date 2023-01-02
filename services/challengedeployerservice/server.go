package challengedeployerservice

import (
	"log"

	g "github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
)

func DeployToAll(localFilePath string, pathInPod string, ns ...string) error {

	if err := GetClient(g.KatanaConfig.KubeConfig); err != nil {
		return err
	}

	// Get pods
	pods, err := getPods(map[string]string{"app": g.ClusterConfig.TeamLabel})
	if err != nil {
		log.Println(err)
		return err
	}

	// Loop over pods
	for _, pod := range pods {
		// Copy file into pod
		if err := utils.CopyIntoPod(pod.Name, g.TeamVmConfig.ContainerName, pathInPod, localFilePath, ns...); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
