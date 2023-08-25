package harbor

import (
	"log"

	"github.com/sdslabs/katana/lib/utils"
)

func SetupHarbor() error {
	client, _ := utils.GetKubeClient()

	deploymentNames := []string{
		"katana-release-harbor-core",
		"katana-release-harbor-portal",
		"katana-release-harbor-registry",
		"katana-release-harbor-jobservice",
	}
	namespace := "katana"

	for _, deploymentName := range deploymentNames {
		if err := utils.WaitForDeploymentReady(client, deploymentName, namespace); err != nil {
			log.Printf("Error testing deployment '%s': %v\n", deploymentName, err)
		}
	}

	if err := updateOrAddHost(); err != nil {
		return err
	}
	if err := setAdminPassword(); err != nil {
		return err
	}
	if err := createHarborProject("katana"); err != nil {
		return err
	}
	if err := setCertificateToDocker(); err != nil {
		return err
	}
	if err := deployHarborClusterDaemonSet(); err != nil {
		return err
	}

	return nil
}
