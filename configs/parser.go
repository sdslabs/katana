package configs

import (
	"flag"
	"os"

	"github.com/BurntSushi/toml"
)

func getConfiguration() *KatanaCfg {
	flag.Parse()
	config := &KatanaCfg{}
	if _, err := toml.DecodeFile(*configFile, config); err != nil {
		os.Exit(1)
	}
	return config
}

var (
	configFile = flag.String("conf", "config.toml", "location of config file")

	KatanaConfig = getConfiguration()

	ServicesConfig = KatanaConfig.Services

	APIConfig = ServicesConfig.API

	VMDeployerServiceConfig = ServicesConfig.VMDeployer
)
