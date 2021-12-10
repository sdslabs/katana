package configs

type API struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type ClusterCfg struct {
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
	SSHProvider       SSHProviderCfg       `toml:"sshprovider"`
}

type TeamChallengeConfig struct {
	ChallengeDir string `toml:"challengedir"`
	TempDir      string `toml:"tmpdir"`
	InitFile     string `toml:"initfile"`
	DaemonPort   uint   `toml:"daemonport"`
}

type SSHProviderCfg struct {
	Host        string `toml:"host"`
	Port        uint   `toml:"port"`
	CredsFile   string `toml:"creds_file"`
	PasswordLen uint   `toml:"password_length"`
}

type MongoCfg struct {
	URL string `toml:"url"`
}

type CloudCfg struct {
	Provider string `toml:"provider"`
}

type AzureCfg struct {
	ResourceGroupName string `toml:"resource_group_name"`
	ClusterName       string `toml:"cluster_name"`
	Location          string `toml:"location"`
}

type GCPCfg struct {
	ProjectID       string `toml:"project_id"`
	CredentialsFile string `toml:"credentials_file"`
	ClusterName     string `toml:"cluster_name"`
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
	CloudConfig   CloudCfg            `toml:"cloud"`
	AzureConfig   AzureCfg            `toml:"azure"`
	GCPConfig     GCPCfg              `toml:"gcp"`
}
