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
	Token string `json:"token"`
}

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	userService := services.NewUserService(global.DB)
	token, err := userService.Login(req.Username, req.Password)
	if err != nil {
		if code.IsErrCode(err) {
			response.HandleError(c, err)
		} else {
			response.InternalServerError(c, "登录失败")
		}
		return
	}

	response.Success(c, &LoginResponse{
		Token: token,
	})
}
