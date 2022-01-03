package main

import (
	"github.com/sdslabs/katana/services/challengedeployerservice"
	"github.com/gin-gonic/gin"
)

func main() {

	// Deploy pods here

	// init HTTP Server
	r := gin.Default()
	// Register challenge deployer routes
	challengedeployerservice.ChallengeDeployerInit(r)
	r.Run()
}
