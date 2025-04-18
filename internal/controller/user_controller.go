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
	// 记录请求开始
	logger.Info("UserLogin 开始处理登录请求",
		"client_ip", ctx.ClientIP(),
		"user_agent", ctx.Request.UserAgent())

	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 记录参数验证错误
		logger.Warn("UserLogin 请求参数验证失败",
			"error", err,
			"client_ip", ctx.ClientIP())
		response.Error(ctx, errorx.ErrInvalidParams.WithError(err))
		return
	}

	// 获取客户端IP地址
	clientIP := ctx.ClientIP()

	// 记录尝试登录信息
	logger.Info("UserLogin 尝试登录",
		"username", req.Username,
		"ip", clientIP)

	tokenResponse, err := ctrl.userService.Login(req.Username, req.Password, clientIP)
	if err != nil {
		// 根据错误类型记录不同级别的日志
		if errorx.IsAppErr(err) {
			appErr := err.(*errorx.AppError)
			errCode := appErr.GetErrorCode()

			// 如果是用户不存在或密码错误，记录为警告
			if errCode == errorx.ErrorUserNotFoundCode || errCode == errorx.ErrorUserPasswordErrorCode {
				logger.Warn("UserLogin 登录失败",
					"error", err,
					"username", req.Username,
					"ip", clientIP,
					"error_code", errCode)
			} else {
				// 其他错误记录为错误
				logger.Error("UserLogin 登录失败",
					"error", err,
					"username", req.Username,
					"ip", clientIP,
					"error_code", errCode)
			}
		} else {
			// 非AppError类型的错误记录为错误
			logger.Error("UserLogin 登录失败",
				"error", err,
				"username", req.Username,
				"ip", clientIP)
		}

		// 返回错误响应
		response.Error(ctx, err)
		return
	}

	// 记录登录成功
	logger.Info("UserLogin 登录成功",
		"username", req.Username,
		"access_token", tokenResponse.AccessToken[:10]+"...", // 只显示令牌前10个字符
		"ip", clientIP)

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
			response.Error(c, errorx.ErrDatabaseQuery)
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
