package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/services"
)

// 获取系统设置

func NewSystemController(systemService *services.SystemService) *SystemController {
	return &SystemController{
		systemService: systemService,
	}
}

type SystemController struct {
	systemService *services.SystemService
}

func (sc *SystemController) GetSystemSettings(c *gin.Context) {
	// 获取系统设置
	settings := sc.systemService.GetSystemSettings(c.Request.Context())

	response.Success(c, settings)
}
