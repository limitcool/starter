package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/pkg/response"
)

func Ping(c *gin.Context) {
	response.Success(c, "pong")
}
