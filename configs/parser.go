package configs

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
)

func getConfiguration() *KatanaCfg {
	flag.Parse()
	config := &KatanaCfg{}
	if _, err := toml.DecodeFile(*configFile, config); err != nil {
		log.Fatal(err)
	}
	return config
}

var (
	configFile = flag.String("conf", "config.toml", "location of config file")

	KatanaConfig = getConfiguration()

	ServicesConfig = KatanaConfig.Services

	APIConfig = ServicesConfig.API

	ClusterConfig = KatanaConfig.Cluster

	ChallengeDeployerConfig = ServicesConfig.ChallengeDeployer

	SSHProviderConfig = ServicesConfig.SSHProvider

	AdminConfig = KatanaConfig.AdminConfig

	TeamVmConfig = KatanaConfig.TeamVmConfig

	MongoConfig = KatanaConfig.Mongo

	CloudConfig = KatanaConfig.CloudConfig
)
