package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/services"
	"github.com/spf13/cast"
)

func NewUserController(sysUserService *services.SysUserService, userService *services.UserService) *UserController {
	return &UserController{
		sysUserService: sysUserService,
		userService:    userService,
	}
}

type UserController struct {
	sysUserService *services.SysUserService
	userService    *services.UserService
}

// UserLogin 普通用户登录
func (ctrl *UserController) UserLogin(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	// 获取客户端IP地址
	clientIP := ctx.ClientIP()
	tokenResponse, err := ctrl.userService.Login(req.Username, req.Password, clientIP)
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
func (ctrl *UserController) UserRegister(c *gin.Context) {
	var req v1.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取客户端IP地址
	clientIP := c.ClientIP()

	user, err := ctrl.userService.Register(req, clientIP)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 隐藏密码等敏感信息
	user.Password = ""

	response.Success(c, user)
}

// UserChangePassword 修改密码
func (ctrl *UserController) UserChangePassword(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("user_id")

	var req v1.UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	err := ctrl.userService.ChangePassword(cast.ToInt64(userID), req.OldPassword, req.NewPassword)
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
func (ctrl *UserController) UserInfo(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errorx.ErrUserNotFound)
		return
	}

	// 获取用户信息
	user, err := ctrl.userService.GetUserByID(cast.ToInt64(userID))
	if err != nil {
		response.Error(c, err)
		return
	}

	// 隐藏密码等敏感信息
	user.Password = ""

	response.Success(c, user)
}

// RefreshToken 刷新访问令牌
func (ctrl *UserController) RefreshToken(c *gin.Context) {
	var req v1.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 使用控制器中的服务实例
	tokenResponse, err := ctrl.sysUserService.RefreshToken(req.RefreshToken)
	if err != nil {
		logger.LogError("RefreshToken 刷新访问令牌失败", err,
			"refresh_token", req.RefreshToken)
		response.Error(c, err)
		return
	}

	response.Success(c, tokenResponse)
}
