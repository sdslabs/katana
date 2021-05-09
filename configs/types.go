package configs

type API struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type Cluster struct {
	BroadcastCount uint   `toml:"broadcastcount"`
	BroadcastLabel string `toml:"broadcastlabel"`
	TeamCount      uint   `toml:"teamcount"`
	TeamLabel      string `toml:"teamlabel"`
}

type ChallengeDeployerConfig struct {
	Host           string `toml:"host"`
	Port           uint   `toml:"port"`
	BroadcastPort  uint   `toml:"broadcastport"`
	TeamClientPort uint   `toml:"teamclientport"`
}

type Services struct {
	API               API                     `toml:"api"`
	ChallengeDeployer ChallengeDeployerConfig `toml:"challengedeployer"`
}

type KatanaCfg struct {
	KubeHost      string   `toml:"kubehost"`
	KubeNameSpace string   `toml:"kubenamesapce"`
	KubeConfig    string   `toml:"kubeconfig"`
	Services      Services `toml:"services"`
	Cluster       Cluster  `toml:"cluster"`
}
