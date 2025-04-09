package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/services"
)

var UserControllerInstance = UserController{
	userService: services.Instance().GetUserService(),
}

type UserController struct {
	userService *services.UserService
}

// UserLogin 普通用户登录
func (uc *UserController) UserLogin(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ParamError(ctx, "无效的请求参数")
		return
	}

	// 获取客户端IP地址
	clientIP := ctx.ClientIP()

	db := services.Instance().GetDB()
	userService := services.NewNormalUserService(db)
	tokenResponse, err := userService.Login(req.Username, req.Password, clientIP)
	if err != nil {
		if errorx.IsErrCode(err) {
			response.HandleError(ctx, err)
		} else {
			response.ServerError(ctx)
		}
		return
	}

	response.Success(ctx, tokenResponse)
}

// UserRegister 普通用户注册
func (uc *UserController) UserRegister(c *gin.Context) {
	var req v1.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "无效的请求参数")
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

	db := services.Instance().GetDB()
	userService := services.NewNormalUserService(db)
	user, err := userService.Register(registerReq)
	if err != nil {
		if errorx.IsErrCode(err) {
			response.HandleError(c, err)
		} else {
			response.ServerError(c)
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
		response.ParamError(c, "无效的请求参数")
		return
	}

	db := services.Instance().GetDB()
	userService := services.NewNormalUserService(db)
	err := userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		if errorx.IsErrCode(err) {
			response.HandleError(c, err)
		} else {
			response.ServerError(c)
		}
		return
	}

	response.Success[any](c, nil)
}

// UserInfo 获取用户信息
func (uc *UserController) UserInfo(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("user_id")

	db := services.Instance().GetDB()
	userService := services.NewNormalUserService(db)
	user, err := userService.GetUserByID(userID.(uint))
	if err != nil {
		if errorx.IsErrCode(err) {
			response.HandleError(c, err)
		} else {
			response.ServerError(c)
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
		response.ParamError(c, "无效的请求参数")
		return
	}

	// 使用服务管理器获取用户服务
	userService := services.Instance().GetUserService()
	tokenResponse, err := userService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, tokenResponse)
}
