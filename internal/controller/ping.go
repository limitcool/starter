package controller

import (
	"github.com/gin-gonic/gin"
)

// Ping 测试ping
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
