package controller

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/services"
	"github.com/spf13/cast"
	"go.uber.org/fx"
)

func NewUserController(params ControllerParams) *UserController {
	controller := &UserController{
		sysUserService: params.SysUserService,
		userService:    params.UserService,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "UserController initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "UserController stopped")
			return nil
		},
	})

	return controller
}

type UserController struct {
	sysUserService *services.SysUserService
	userService    *services.UserService
}

// UserLogin 普通用户登录
func (ctrl *UserController) UserLogin(ctx *gin.Context) {
	// 获取请求上下文
	reqCtx := ctx.Request.Context()

	// 记录请求开始
	logger.InfoContext(reqCtx, "UserLogin 开始处理登录请求",
		"client_ip", ctx.ClientIP(),
		"user_agent", ctx.Request.UserAgent())

	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 记录参数验证错误
		logger.WarnContext(reqCtx, "UserLogin 请求参数验证失败",
			"error", err,
			"client_ip", ctx.ClientIP())
		response.Error(ctx, errorx.ErrInvalidParams.WithError(err))
		return
	}

	// 获取客户端IP地址
	clientIP := ctx.ClientIP()

	// 记录尝试登录信息
	logger.InfoContext(reqCtx, "UserLogin 尝试登录",
		"username", req.Username,
		"ip", clientIP)

	tokenResponse, err := ctrl.userService.Login(ctx.Request.Context(), req.Username, req.Password, clientIP)
	if err != nil {
		// 根据错误类型记录不同级别的日志
		if errorx.IsAppErr(err) {
			appErr := err.(*errorx.AppError)
			errCode := appErr.GetErrorCode()

			// 如果是用户不存在或密码错误，记录为警告
			if errCode == errorx.ErrorUserNotFoundCode || errCode == errorx.ErrorUserPasswordErrorCode {
				logger.WarnContext(reqCtx, "UserLogin 登录失败",
					"error", err,
					"username", req.Username,
					"ip", clientIP,
					"error_code", errCode)
			} else {
				// 其他错误记录为错误
				logger.ErrorContext(reqCtx, "UserLogin 登录失败",
					"error", err,
					"username", req.Username,
					"ip", clientIP,
					"error_code", errCode)
			}
		} else {
			// 非AppError类型的错误记录为错误
			logger.ErrorContext(reqCtx, "UserLogin 登录失败",
				"error", err,
				"username", req.Username,
				"ip", clientIP)
		}

		// 返回错误响应
		response.Error(ctx, err)
		return
	}

	// 记录登录成功
	logger.InfoContext(reqCtx, "UserLogin 登录成功",
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

	user, err := ctrl.userService.Register(c.Request.Context(), req, clientIP)
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

	err := ctrl.userService.ChangePassword(c.Request.Context(), cast.ToInt64(userID), req.OldPassword, req.NewPassword)
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
	user, err := ctrl.userService.GetUserByID(c.Request.Context(), cast.ToInt64(userID))
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
	tokenResponse, err := ctrl.sysUserService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		logger.LogErrorContext(c.Request.Context(), "RefreshToken 刷新访问令牌失败", err,
			"refresh_token", req.RefreshToken)
		response.Error(c, err)
		return
	}

	response.Success(c, tokenResponse)
}
