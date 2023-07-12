package harbor

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

	if err := dockerLogin(); err != nil {
		return err
	}

	if err := setHostsInCluster(); err != nil {
		return err
	}

	return nil
}
