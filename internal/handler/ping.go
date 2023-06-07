package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/pkg/code"
)

// Ping ping
// @Summary ping
// @Description ping
// @Tags system
// @Accept  json
// @Produce  json
// @Router /ping [get]
func Ping(c *gin.Context) {
	code.AutoResponse(c, "pong", nil)
}
