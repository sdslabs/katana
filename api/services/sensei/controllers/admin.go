package controllers

import (
	"github.com/gin-gonic/gin"
)

func HelloAdmin(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
		"data":    "hello admin",
	})
}
