package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/pkg/code"
)

func Ping(c *gin.Context) {
	code.AutoResponse(c, "pong", nil)
}
