package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/services"
)

// AdminController 管理系统控制器
type AdminController struct {
	adminService *services.AdminService
}

// NewAdminController 创建管理系统控制器
func NewAdminController(adminService *services.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

// GetSystemSettings 获取系统设置
func (c *AdminController) GetSystemSettings(ctx *gin.Context) {
	// 获取系统设置
	settings := c.adminService.GetSystemSettings(ctx.Request.Context())

	response.Success(ctx, settings)
}
