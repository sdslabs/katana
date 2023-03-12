package configs

import (
	"flag"
	"log"
	"os"

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

func WriteConfiguration() {
	flag.Parse()
	config := KatanaConfig
	file, err := os.OpenFile(*configFile, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if err := toml.NewEncoder(file).Encode(config); err != nil {
		log.Fatal(err)
	}
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

	MySQLConfig = KatanaConfig.MySQL
)
