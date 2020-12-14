package configs

type API struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type VMDeployerService struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type ChallengeDeployerConfig struct {
	Host           string `toml:"host"`
	Port           int    `toml:"port"`
	TeamLabel      string `toml:"teamlabel"`
	BroadcastPort  int    `toml:"broadcastport"`
	TeamClientPort int    `toml:"teamclientport"`
	KubeHost       string `toml:"kubehost"`
	KubeNameSpace  string `toml:"kubenamesapce"`
	KubeConfig     string `toml:"kubeconfig"`
}

type Services struct {
	API               API                     `toml:"api"`
	VMDeployer        VMDeployerService       `toml:"vmdeployer"`
	ChallengeDeployer ChallengeDeployerConfig `toml:"challengedeployer"`
}

type KatanaCfg struct {
	Services Services `toml:"services"`
}
