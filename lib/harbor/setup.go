package harbor

func SetupHarbor() error {
	if !checkHarborHostsEntryExists() {
		if err := addHarborHostsEntry(); err != nil {
			return err
		}
	}

	if err := createHarborProject("katana"); err != nil {
		return err
	}

	if err := getHarborCertificate(); err != nil {
		return err
	}

	if err := addCertificateToDocker(); err != nil {
		return err
	}

	if err := setupDockerCredentials(); err != nil {
		return err
	}

	return nil
}
