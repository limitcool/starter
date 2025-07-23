package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/errorx"
	"github.com/limitcool/starter/internal/middleware"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// UserHandler 用户处理器
type UserHandler struct {
	*BaseHandler
	app         AppContext
	authService *AuthService
}

var _ RouterInitializer = (*UserHandler)(nil) // 用于接口断言，_ 变量编译后会被移除

// NewUserHandler 创建用户处理器
func NewUserHandler(app AppContext) *UserHandler {
	handler := &UserHandler{
		BaseHandler: NewBaseHandler(app.GetDB(), app.GetConfig()),
		authService: NewAuthService(app.GetConfig()), // TODO: service 应该移到 services 文件夹
		app:         app,
	}

	handler.LogInit("UserHandler")
	return handler
}

func (h *UserHandler) InitRouters(g *gin.RouterGroup, root *gin.Engine) {

	// 公共路由
	public := g.Group("")
	{
		// 用户登录（管理员和普通用户使用同一接口）
		public.POST("/login", h.UserLogin)

		// 用户注册
		public.POST("/register", h.UserRegister)
	}

	// 需要认证的路由
	authenticated := g.Group("", middleware.JWTAuth(h.Config))

	// 普通用户路由 - 使用JWT认证
	user := authenticated.Group("/user")
	{
		// 用户信息
		user.GET("/info", h.UserInfo)

		// 修改密码
		user.POST("/change-password", h.UserChangePassword)
	}
}

// UserLogin 用户登录
func (h *UserHandler) UserLogin(ctx *gin.Context) {
	// 获取请求上下文
	reqCtx := ctx.Request.Context()

	// 记录请求开始
	logger.InfoContext(reqCtx, "UserLogin request started",
		"client_ip", ctx.ClientIP(),
		"user_agent", ctx.Request.UserAgent())

	var req dto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 记录参数验证错误
		logger.WarnContext(reqCtx, "UserLogin request validation failed",
			"error", err,
			"client_ip", ctx.ClientIP())
		response.Error(ctx, errorx.ErrInvalidParams.New(ctx, struct{ Params string }{err.Error()}))
		return
	}

	// 获取客户端IP地址
	clientIP := ctx.ClientIP()

	// 记录尝试登录信息
	logger.InfoContext(reqCtx, "UserLogin attempting login",
		"username", req.Username,
		"ip", clientIP)

	// 创建用户仓库
	userRepo := model.NewUserRepo(h.DB)

	// 查询用户
	user, err := userRepo.GetByUsername(reqCtx, req.Username)
	if err != nil {
		if errorx.ErrUserNotFound.Is(err) {
			// 用户不存在
			logger.WarnContext(reqCtx, "UserLogin user not found",
				"username", req.Username,
				"ip", clientIP)
			response.Error(ctx, err)
			return
		}
		// 数据库错误
		logger.ErrorContext(reqCtx, "UserLogin failed to query user",
			"error", err,
			"username", req.Username,
			"ip", clientIP)
		response.Error(ctx, err)
		return
	}

	// 检查用户是否启用
	if !user.Enabled {
		disabledErr := errorx.ErrUserDisabled.New(ctx, struct{ Name string }{req.Username})
		logger.WarnContext(reqCtx, "UserLogin user is disabled",
			"username", req.Username,
			"ip", clientIP)
		response.Error(ctx, disabledErr)
		return
	}

	// 验证密码
	if !crypto.CheckPassword(user.Password, req.Password) {
		passwordErr := errorx.ErrUserPassword.New(ctx, struct{ Name string }{req.Username})
		logger.WarnContext(reqCtx, "UserLogin password incorrect",
			"username", req.Username,
			"ip", clientIP)
		response.Error(ctx, passwordErr)
		return
	}

	// 更新最后登录时间和IP
	if err := userRepo.UpdateLastLogin(reqCtx, int64(user.ID), clientIP); err != nil {
		logger.WarnContext(reqCtx, "UserLogin failed to update login info",
			"error", err,
			"username", req.Username,
			"ip", clientIP)
		// 这里不返回错误，因为登录信息更新失败不应该影响用户登录
	}

	// 获取用户角色
	var roles []string
	if user.IsAdmin {
		roles = []string{"admin"}
	} else {
		roles = []string{"user"}
	}

	// 生成令牌
	tokenResponse, err := h.authService.GenerateTokensWithContext(reqCtx, uint(user.ID), user.Username, user.IsAdmin, roles)
	if err != nil {
		logger.ErrorContext(reqCtx, "UserLogin failed to generate token",
			"error", err,
			"username", req.Username,
			"ip", clientIP)
		response.Error(ctx, errorx.ErrGenVisitToken.New(ctx, errorx.None))
		return
	}

	// 记录登录成功
	logger.InfoContext(reqCtx, "UserLogin successful",
		"username", req.Username,
		"access_token", tokenResponse.AccessToken[:10]+"...", // 只显示令牌前10个字符
		"ip", clientIP)

	response.Success(ctx, tokenResponse)
}

// UserRegister 用户注册
func (h *UserHandler) UserRegister(ctx *gin.Context) {
	// 获取请求上下文
	reqCtx := ctx.Request.Context()

	var req dto.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.WarnContext(reqCtx, "UserRegister request validation failed",
			"error", err,
			"client_ip", ctx.ClientIP())
		response.Error(ctx, errorx.ErrInvalidParams.New(ctx, struct{ Params string }{err.Error()}))
		return
	}

	// 获取客户端IP地址
	clientIP := ctx.ClientIP()

	// 创建用户仓库
	userRepo := model.NewUserRepo(h.DB)

	// 检查用户名是否已存在
	exists, err := userRepo.IsExist(reqCtx, req.Username)
	if err != nil {
		logger.ErrorContext(reqCtx, "UserRegister failed to check username existence",
			"error", err,
			"username", req.Username,
			"ip", clientIP)
		response.Error(ctx, err)
		return
	}

	if exists {
		existsErr := errorx.ErrUserExists.New(ctx, struct{ Name string }{req.Username})
		logger.WarnContext(reqCtx, "UserRegister username already exists",
			"username", req.Username,
			"ip", clientIP)
		response.Error(ctx, existsErr)
		return
	}

	// 哈希密码
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		logger.ErrorContext(reqCtx, "UserRegister failed to hash password",
			"error", err,
			"username", req.Username,
			"ip", clientIP)
		response.Error(ctx, errorx.ErrPasswordEncrypt.New(ctx, errorx.None).Wrap(err))
		return
	}

	// 创建用户
	user := &model.User{
		Username:   req.Username,
		Password:   hashedPassword,
		Nickname:   req.Nickname,
		Email:      req.Email,
		Mobile:     req.Mobile,
		Enabled:    true,
		Gender:     req.Gender,
		Birthday:   &time.Time{},
		Address:    req.Address,
		RegisterIP: clientIP,
		IsAdmin:    false, // 普通用户注册，不是管理员
	}

	if err := userRepo.Create(reqCtx, user); err != nil {
		logger.ErrorContext(reqCtx, "UserRegister failed to create user",
			"error", err,
			"username", req.Username,
			"ip", clientIP)
		response.Error(ctx, err)
		return
	}

	// 隐藏密码等敏感信息
	user.Password = ""

	logger.InfoContext(reqCtx, "UserRegister user registration successful",
		"username", req.Username,
		"ip", clientIP)

	response.Success(ctx, user)
}

// UserInfo 获取用户信息
func (h *UserHandler) UserInfo(ctx *gin.Context) {
	// 获取用户ID
	id, ok := h.Helper.GetUserID(ctx)
	if !ok {
		return
	}

	// 创建用户仓库
	userRepo := model.NewUserRepo(h.DB)

	// 查询用户信息
	user, err := userRepo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errorx.ErrUserNotFound.Is(err) {
			h.Helper.HandleNotFoundError(ctx, err, "UserInfo", "user_id", id)
			return
		}
		h.Helper.HandleDBError(ctx, err, "UserInfo", "user_id", id)
		return
	}

	// 隐藏敏感信息
	user.Password = ""

	h.Helper.LogSuccess(ctx, "UserInfo", "user_id", id)
	response.Success(ctx, user)
}

// UserChangePassword 修改密码
func (h *UserHandler) UserChangePassword(ctx *gin.Context) {
	// 获取用户ID
	id, ok := h.Helper.GetUserID(ctx)
	if !ok {
		return
	}

	// 绑定请求参数
	var req dto.UserChangePasswordRequest
	if !h.Helper.BindJSON(ctx, &req, "UserChangePassword") {
		return
	}

	// 创建用户仓库
	userRepo := model.NewUserRepo(h.DB)

	// 查询用户
	user, err := userRepo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errorx.ErrUserNotFound.Is(err) {
			h.Helper.HandleNotFoundError(ctx, err, "UserChangePassword", "user_id", id)
			return
		}
		h.Helper.HandleDBError(ctx, err, "UserChangePassword", "user_id", id)
		return
	}

	// 验证旧密码
	if !crypto.CheckPassword(user.Password, req.OldPassword) {
		h.Helper.LogWarning(ctx, "UserChangePassword old password incorrect", "user_id", id)
		response.Error(ctx, errorx.ErrOldPasswordError.New(ctx, errorx.None))
		return
	}

	// 哈希新密码
	hashedPassword, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		h.Helper.LogError(ctx, "UserChangePassword failed to hash password", "error", err, "user_id", id)
		response.Error(ctx, errorx.ErrPasswordEncrypt.New(ctx, errorx.None))
		return
	}

	// 更新密码
	if err := userRepo.UpdatePassword(ctx.Request.Context(), id, hashedPassword); err != nil {
		h.Helper.HandleDBError(ctx, err, "UserChangePassword", "user_id", id)
		return
	}

	h.Helper.LogSuccess(ctx, "UserChangePassword", "user_id", id)
	response.SuccessNoData(ctx, "密码修改成功")
}
