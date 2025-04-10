package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/services"
)

var AdminControllerInstance AdminController = AdminController{}

type AdminController struct {
}

// AdminLogin 管理员登录
func (ac *AdminController) AdminLogin(c *gin.Context) {
	var req v1.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError("AdminLogin 无效的请求参数", err)
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取客户端IP地址
	clientIP := c.ClientIP()

	// 使用服务管理器获取用户服务
	userService := services.NewSysUserService()
	tokenResponse, err := userService.Login(req.Username, req.Password, clientIP)
	if err != nil {
		// 使用辅助函数记录错误，同时包含额外的上下文信息
		logger.LogError("AdminLogin 登录失败", err,
			"username", req.Username,
			"ip", clientIP)

		// 直接返回包装后的错误
		response.Error(c, err)
		return
	}

	response.Success(c, tokenResponse)
}
