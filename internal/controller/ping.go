package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/pkg/apiresponse"
)

// Ping 健康检查
func Ping(c *gin.Context) {
	apiresponse.Success(c, "pong")
}
