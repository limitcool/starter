package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/pkg/code"
	"github.com/limitcool/starter/internal/services"
)

// UserRegisterRequest 用户注册请求参数
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Gender   string `json:"gender"`
	Address  string `json:"address"`
}

// UserLogin 普通用户登录
func UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "无效的请求参数")
		return
	}

	// 获取客户端IP地址
	clientIP := c.ClientIP()

	db := services.Instance().GetDB()
	userService := services.NewNormalUserService(db)
	tokenResponse, err := userService.Login(req.Username, req.Password, clientIP)
	if err != nil {
		if code.IsErrCode(err) {
			response.HandleError(c, err)
		} else {
			response.ServerError(c)
		}
		return
	}

	response.Success(c, tokenResponse)
}

// UserRegister 普通用户注册
func UserRegister(c *gin.Context) {
	var req UserRegisterRequest
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
		if code.IsErrCode(err) {
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
func UserChangePassword(c *gin.Context) {
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
		if code.IsErrCode(err) {
			response.HandleError(c, err)
		} else {
			response.ServerError(c)
		}
		return
	}

	response.Success[any](c, nil)
}

// UserInfo 获取用户信息
func UserInfo(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("user_id")

	db := services.Instance().GetDB()
	userService := services.NewNormalUserService(db)
	user, err := userService.GetUserByID(userID.(uint))
	if err != nil {
		if code.IsErrCode(err) {
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
