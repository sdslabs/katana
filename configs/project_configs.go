package configs

import "flag"

var (
	// configFile is the main configuration file for katana
	configFile = flag.String("conf", "config.toml", "location of config file")

	// KatanaConfig is parsed data for `configFile`
	KatanaConfig = getConfiguration()

	// VMDeployerConfig is the configuration for VMDeployer
	VMDeployerConfig = KatanaConfig.VMDeployer
)
