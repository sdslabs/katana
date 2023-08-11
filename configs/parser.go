package configs

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
)

func LoadConfiguration() {
	flag.Parse()
	config := &KatanaCfg{}
	configFile := flag.String("conf", "config.toml", "location of config file")
	if _, err := toml.DecodeFile(*configFile, config); err != nil {
		log.Fatal(err)
	}
	KatanaConfig = config
	ServicesConfig = &KatanaConfig.Services
	APIConfig = &ServicesConfig.API
	ClusterConfig = &KatanaConfig.Cluster
	SSHProviderConfig = &ServicesConfig.SSHProvider
	AdminConfig = &KatanaConfig.AdminConfig
	TeamVmConfig = &KatanaConfig.TeamVmConfig
	MongoConfig = &KatanaConfig.Mongo
	MySQLConfig = &KatanaConfig.MySQL
	HarborConfig = &KatanaConfig.Harbor
}

var (
	KatanaConfig *KatanaCfg

	ServicesConfig *ServicesCfg

	APIConfig *API

	ClusterConfig *ClusterCfg

	SSHProviderConfig *SSHProviderCfg

	AdminConfig *AdminCfg

	TeamVmConfig *TeamChallengeConfig

	MongoConfig *MongoCfg

	MySQLConfig *MySQLCfg

	HarborConfig *HarborCfg
)
