package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/services"
	"github.com/spf13/cast"
)

func NewUserController(userService *services.SysUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

type UserController struct {
	userService *services.SysUserService
}

// UserLogin 普通用户登录
func (uc *UserController) UserLogin(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	// 获取客户端IP地址
	clientIP := ctx.ClientIP()
	userService := services.NewUserService()
	tokenResponse, err := userService.Login(req.Username, req.Password, clientIP)
	if err != nil {
		logger.LogError("UserLogin 登录失败", err,
			"username", req.Username,
			"ip", clientIP)
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, tokenResponse)
}

// UserRegister 普通用户注册
func (uc *UserController) UserRegister(c *gin.Context) {
	var req v1.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取客户端IP地址
	clientIP := c.ClientIP()

	userService := services.NewUserService()
	user, err := userService.Register(req, clientIP)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 隐藏密码等敏感信息
	user.Password = ""

	response.Success(c, user)
}

// UserChangePassword 修改密码
func (uc *UserController) UserChangePassword(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("user_id")

	var req v1.UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	userService := services.NewUserService()
	err := userService.ChangePassword(cast.ToInt64(userID), req.OldPassword, req.NewPassword)
	if err != nil {
		if errorx.IsAppErr(err) {
			response.Error(c, err)
		} else {
			response.Error(c, errorx.ErrDatabaseQueryError)
		}
		return
	}

	response.Success[any](c, nil)
}

// UserInfo 获取用户信息
func (uc *UserController) UserInfo(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errorx.ErrUserNotFound)
		return
	}

	user, err := model.NewUser().GetUserByID(cast.ToInt64(userID))
	if err != nil {
		response.Error(c, err)
		return
	}

	// 隐藏密码等敏感信息
	user.Password = ""

	response.Success(c, user)
}

// RefreshToken 刷新访问令牌
func (uc *UserController) RefreshToken(c *gin.Context) {
	var req v1.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 使用控制器中的服务实例
	tokenResponse, err := uc.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		logger.LogError("RefreshToken 刷新访问令牌失败", err,
			"refresh_token", req.RefreshToken)
		response.Error(c, err)
		return
	}

	response.Success(c, tokenResponse)
}
