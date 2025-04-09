package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/services"
)

var AdminControllerInstance AdminController = AdminController{}

type AdminController struct {
}

// AdminLogin 管理员登录
func (ac *AdminController) AdminLogin(c *gin.Context) {
	var req v1.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "无效的请求参数")
		return
	}

	// 获取客户端IP地址
	clientIP := c.ClientIP()

	// 使用服务管理器获取用户服务
	userService := services.Instance().GetUserService()
	tokenResponse, err := userService.Login(req.Username, req.Password, clientIP)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, tokenResponse)
}
