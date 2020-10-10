package configs

//VMDeployer is the configuration for VmDeployer microservice
type VMDeployer struct {
	Jobname         string `toml:"jobname"`
	Region          string `toml:"region"`
	Priority        int    `toml:"priority"`
	Datacentre      string `toml:"datacentre"`
	GroupName       string `toml:"groupName"`
	RestartAttempts int    `toml:"restartAttempts"`
	RestartMode     string `toml:"restartMode"`
	Taskname        string `toml:"taskname"`
	Driver          string `toml:"driver"`
	NetworkName     string `toml:"networkName"`
	FileNum         int    `toml:"fileNum"`
	FileSize        int    `toml:"fileSize"`
}

// KatanaCfg is the configuration for the entire project
type KatanaCfg struct {
	VMDeployer VMDeployer `toml:"vmdeployer"`
}
