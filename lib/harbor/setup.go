package harbor

import (
	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
)

func SetupHarbor() error {
	if !checkHarborHostsEntryExists() {
		if err := addHarborHostsEntry(); err != nil {
			return err
		}
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

	if err := utils.DockerLogin(configs.KatanaConfig.Harbor.Username, configs.KatanaConfig.Harbor.Password); err != nil {
		return err
	}

	if err := setHostsInCluster(); err != nil {
		return err
	}

	return nil
}