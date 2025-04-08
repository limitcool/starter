package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/pkg/code"
	"github.com/limitcool/starter/pkg/response"
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
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 获取客户端IP地址
	clientIP := c.ClientIP()

	userService := services.NewUserService(global.DB)
	tokenResponse, err := userService.Login(req.Username, req.Password, clientIP)
	if err != nil {
		if code.IsErrCode(err) {
			response.HandleError(c, err)
		} else {
			response.InternalServerError(c, "登录失败")
		}
		return
	}

	response.Success(c, tokenResponse)
}

// RefreshToken 刷新访问令牌
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	userService := services.NewUserService(global.DB)
	tokenResponse, err := userService.RefreshToken(req.RefreshToken)
	if err != nil {
		if code.IsErrCode(err) {
			response.HandleError(c, err)
		} else {
			response.InternalServerError(c, "刷新令牌失败")
		}
		return
	}

	response.Success(c, tokenResponse)
}
