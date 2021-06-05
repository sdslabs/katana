package configs

type API struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type Cluster struct {
	BroadcastCount uint     `toml:"broadcastcount"`
	BroadcastLabel string   `toml:"broadcastlabel"`
	TeamCount      uint     `toml:"teamcount"`
	TeamLabel      string   `toml:"teamlabel"`
	ManifestDir    string   `toml:"manifest_dir"`
	Manifests      []string `toml:"manifests"`
}

type ChallengeDeployerConfig struct {
	Host           string `toml:"host"`
	Port           uint   `toml:"port"`
	BroadcastPort  uint   `toml:"broadcastport"`
	TeamClientPort uint   `toml:"teamclientport"`
	ArtifactLabel  string `toml:"challenge_artifact_label"`
}

type Services struct {
	API               API                     `toml:"api"`
	ChallengeDeployer ChallengeDeployerConfig `toml:"challengedeployer"`
}

type KatanaCfg struct {
	KubeHost      string              `toml:"kubehost"`
	KubeNameSpace string              `toml:"kubenamespace"`
	KubeConfig    string              `toml:"kubeconfig"`
	Services      Services            `toml:"services"`
	Cluster       Cluster             `toml:"cluster"`
	TeamVmConfig  TeamChallengeConfig `toml:"teamvm"`
}

type TeamChallengeConfig struct {
	ChallengeDir string `toml:"challengedir"`
	TempDir      string `toml:"tmpdir"`
	InitFile     string `toml:"initfile"`
	DaemonPort   uint   `toml:"daemonport"`
}
