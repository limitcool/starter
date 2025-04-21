package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/services"
)

// 获取系统设置

func NewAdminSystemController(systemService *services.AdminSystemService) *AdminSystemController {
	return &AdminSystemController{
		systemService: systemService,
	}
}

type AdminSystemController struct {
	systemService *services.AdminSystemService
}

func (sc *AdminSystemController) GetSystemSettings(c *gin.Context) {
	// 获取系统设置
	settings := sc.systemService.GetSystemSettings(c.Request.Context())

	response.Success(c, settings)
}
