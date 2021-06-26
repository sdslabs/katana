package configs

type API struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type Cluster struct {
	DeploymentLabel string   `toml:"deploymentlabel"`
	BroadcastCount  uint     `toml:"broadcastcount"`
	BroadcastLabel  string   `toml:"broadcastlabel"`
	TeamCount       uint     `toml:"teamcount"`
	TeamLabel       string   `toml:"teamlabel"`
	ManifestDir     string   `toml:"manifest_dir"`
	Manifests       []string `toml:"manifests"`
}

type ChallengeDeployerCfg struct {
	Host           string `toml:"host"`
	Port           uint   `toml:"port"`
	BroadcastPort  uint   `toml:"broadcastport"`
	TeamClientPort uint   `toml:"teamclientport"`
	ArtifactLabel  string `toml:"challengeartifactlabel"`
}

type AdminCfg struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type ServicesCfg struct {
	API               API                  `toml:"api"`
	ChallengeDeployer ChallengeDeployerCfg `toml:"challengedeployer"`
}

type TeamChallengeConfig struct {
	ChallengeDir string `toml:"challengedir"`
	TempDir      string `toml:"tmpdir"`
	InitFile     string `toml:"initfile"`
	DaemonPort   uint   `toml:"daemonport"`
}

type KatanaCfg struct {
	KubeHost      string              `toml:"kubehost"`
	KubeNameSpace string              `toml:"kubenamespace"`
	KubeConfig    string              `toml:"kubeconfig"`
	LogFile       string              `toml:"logfile"`
	Services      ServicesCfg         `toml:"services"`
	Cluster       ClusterCfg          `toml:"cluster"`
	Mongo         MongoCfg            `toml:"mongo"`
	TeamVmConfig  TeamChallengeConfig `toml:"teamvm"`
	AdminConfig   AdminCfg            `toml:"admin"`
}
