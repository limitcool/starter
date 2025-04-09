package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/internal/pkg/apiresponse"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应参数
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshTokenRequest 刷新令牌请求参数
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ParamError(c, "无效的请求参数")
		return
	}

	// 获取客户端IP地址
	clientIP := c.ClientIP()

	// 使用服务管理器获取用户服务
	userService := services.Instance().GetUserService()
	tokenResponse, err := userService.Login(req.Username, req.Password, clientIP)
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success(c, tokenResponse)
}

// RefreshToken 刷新访问令牌
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ParamError(c, "无效的请求参数")
		return
	}

	// 使用服务管理器获取用户服务
	userService := services.Instance().GetUserService()
	tokenResponse, err := userService.RefreshToken(req.RefreshToken)
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success(c, tokenResponse)
}
