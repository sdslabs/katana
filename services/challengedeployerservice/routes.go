package challengedeployerservice

import (
	"fmt"

	c "github.com/sdslabs/katana/services/challengedeployerservice/controllers"
	u "github.com/sdslabs/katana/services/challengedeployerservice/utils"
	"github.com/gin-gonic/gin"
)

func ChallengeDeployerInit(router gin.IRouter){
	fmt.Println("Initiating Challenge Deployer")
	if err := u.GetClient(u.KatanaConfig.KubeConfig); err != nil {
		return
	}
    cd := router.Group("/cd")
	cd.Use()
	{
		router.POST("/deploy",c.Deploy)
	}
}


