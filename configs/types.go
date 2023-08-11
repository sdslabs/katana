package configs

type API struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type ClusterCfg struct {
	DeploymentLabel      string   `toml:"deploymentlabel"`
	TeamCount            uint     `toml:"teamcount"`
	TeamLabel            string   `toml:"teamlabel"`
	TemplatedManifestDir string   `toml:"templated_manifest_dir"`
	TemplatedManifests   []string `toml:"templated_manifests"`
}

type AdminCfg struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type ServicesCfg struct {
	API         API            `toml:"api"`
	SSHProvider SSHProviderCfg `toml:"sshprovider"`
}

type TeamChallengeConfig struct {
	TeamPodName   string `toml:"teampodname"`
	ContainerName string `toml:"containername"`
	ChallengeDir  string `toml:"challengedir"`
	TempDir       string `toml:"tmpdir"`
	InitFile      string `toml:"initfile"`
	DaemonPort    uint   `toml:"daemonport"`
}

type SSHProviderCfg struct {
	Host        string `toml:"host"`
	Port        uint   `toml:"port"`
	CredsFile   string `toml:"creds_file"`
	PasswordLen uint   `toml:"password_length"`
}

type MongoCfg struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Port     string `toml:"port"`
	Version  string `toml:"mongosh_version"`
}

type MySQLCfg struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Port     string `toml:"port"`
}

type KatanaCfg struct {
	KubeHost      string              `toml:"kubehost"`
	BackendUrl    string              `toml:"backendurl"`
	KubeNameSpace string              `toml:"kubenamespace"`
	KubeConfig    string              `toml:"kubeconfig"`
	LogFile       string              `toml:"logfile"`
	Services      ServicesCfg         `toml:"services"`
	Cluster       ClusterCfg          `toml:"cluster"`
	Mongo         MongoCfg            `toml:"mongo"`
	TeamVmConfig  TeamChallengeConfig `toml:"teamvm"`
	AdminConfig   AdminCfg            `toml:"admin"`
	MySQL         MySQLCfg            `toml:"mysql"`
	Harbor        HarborCfg           `toml:"harbor"`
	TimeOut       int                 `toml:"timeout"`
}

type HarborCfg struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}
