package controller

import (
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/services"
)

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewSysUserService(),
	}
}

type UserController struct {
	userService *services.SysUserService
}

// UserLogin 普通用户登录
func (uc *UserController) UserLogin(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("UserLogin 无效的请求参数", "err", err)
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	// 获取客户端IP地址
	clientIP := ctx.ClientIP()
	userService := services.NewUserService()
	tokenResponse, err := userService.Login(req.Username, req.Password, clientIP)
	if err != nil {
		if errorx.IsAppErr(err) {
			response.Error(ctx, err)
		} else {
			response.Error(ctx, errorx.ErrDatabaseQueryError)
		}
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

	// 创建注册请求
	registerReq := services.RegisterRequest{
		Username:   req.Username,
		Password:   req.Password,
		Nickname:   req.Nickname,
		Email:      req.Email,
		Mobile:     req.Mobile,
		Gender:     req.Gender,
		Address:    req.Address,
		RegisterIP: clientIP,
	}

	userService := services.NewUserService()
	user, err := userService.Register(registerReq)
	if err != nil {
		if errorx.IsAppErr(err) {
			response.Error(c, err)
		} else {
			response.Error(c, errorx.ErrDatabaseQueryError)
		}
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

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	userService := services.NewUserService()
	err := userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword)
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
	userID, _ := c.Get("user_id")

	userService := services.NewUserService()
	user, err := userService.GetUserByID(userID.(uint))
	if err != nil {
		if errorx.IsAppErr(err) {
			response.Error(c, err)
		} else {
			response.Error(c, errorx.ErrDatabaseQueryError)
		}
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

	// 使用服务管理器获取用户服务
	userService := services.NewSysUserService()
	tokenResponse, err := userService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, tokenResponse)
}
