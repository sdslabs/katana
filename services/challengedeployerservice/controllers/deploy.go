package challengedeployerservice

import (
	"fmt"

	"github.com/gin-gonic/gin"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"

	u "github.com/sdslabs/katana/services/challengedeployerservice/utils"
)

type Params struct {
	USERNAME      string `json:"Username" binding:"required"`
	TOKEN         string `json:"Token" binding:"required"`
	REPOSITORY    string `json:"Repository" binding:"required"`
	CHALLENGENAME string `json:"ChallengeName" binding:"required"`
}

func Deploy(c *gin.Context) {
	var params Params
	c.BindJSON(&params)

	auth := &githttp.BasicAuth{
		Username: params.USERNAME,
		Password: params.TOKEN,
	}

	fmt.Println(params)
	if err := u.Clone(params.REPOSITORY, params.CHALLENGENAME, auth); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": err,
		})
		return
	}

	if err := u.Broadcast(fmt.Sprintf("%s.zip", params.CHALLENGENAME)); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Could not broadcast",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "broadcasted",
	})
}
