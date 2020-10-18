package configs

type API struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type VMDeployerService struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type Services struct {
	API        API               `toml:"api"`
	VMDeployer VMDeployerService `toml:"vmdeployer"`
}

type KatanaCfg struct {
	Services Services `toml:"services"`
}
