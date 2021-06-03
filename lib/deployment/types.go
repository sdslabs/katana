package deployment

// DeploymentConfig contains the total data to be injected in manifest templates
type DeploymentConfig struct {
	TeamCount             uint
	TeamLabel             string
	BroadcastCount        uint
	BroadcastLabel        string
	BroadcastPort         uint
	KubeNameSpace         string
	FluentHost            string
	ChallengeDeployerHost string
}
