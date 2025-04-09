package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
)

// Ping 健康检查
func Ping(c *gin.Context) {
	response.Success(c, "pong")
}
