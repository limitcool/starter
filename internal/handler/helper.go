package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/spf13/cast"
)

// HandlerHelper 处理器辅助工具
type HandlerHelper struct{}

// NewHandlerHelper 创建处理器辅助工具
func NewHandlerHelper() *HandlerHelper {
	return &HandlerHelper{}
}

// GetUserID 从上下文中获取用户ID，如果不存在则返回错误响应
func (h *HandlerHelper) GetUserID(ctx *gin.Context) (int64, bool) {
	reqCtx := ctx.Request.Context()

	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.WarnContext(reqCtx, "用户ID不存在")
		response.Error(ctx, errorx.ErrUserNotLogin.New(ctx, errorx.None))
		return 0, false
	}

	return cast.ToInt64(userID), true
}

// BindJSON 绑定JSON参数，如果失败则返回错误响应
func (h *HandlerHelper) BindJSON(ctx *gin.Context, req interface{}, operation string) bool {
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.WarnContext(reqCtx, operation+" request validation failed",
			"error", err,
			"client_ip", ctx.ClientIP())
		response.Error(ctx, errorx.ErrInvalidParams.New(ctx, struct{ Params string }{err.Error()}))
		return false
	}

	return true
}

// HandleDBError 处理数据库错误，统一日志记录和错误响应
func (h *HandlerHelper) HandleDBError(ctx *gin.Context, err error, operation string, fields ...interface{}) {
	reqCtx := ctx.Request.Context()

	// 构建日志字段
	logFields := []interface{}{
		"error", err,
		"operation", operation,
	}
	logFields = append(logFields, fields...)

	logger.ErrorContext(reqCtx, operation+" database operation failed", logFields...)
	response.Error(ctx, err)
}

// HandleNotFoundError 处理资源不存在错误
func (h *HandlerHelper) HandleNotFoundError(ctx *gin.Context, err error, operation string, fields ...interface{}) {
	reqCtx := ctx.Request.Context()

	// 构建日志字段
	logFields := []interface{}{
		"operation", operation,
	}
	logFields = append(logFields, fields...)

	logger.WarnContext(reqCtx, operation+" resource not found", logFields...)
	response.Error(ctx, err)
}

// LogSuccess 记录成功操作日志
func (h *HandlerHelper) LogSuccess(ctx *gin.Context, operation string, fields ...interface{}) {
	reqCtx := ctx.Request.Context()

	// 构建日志字段
	logFields := []interface{}{
		"operation", operation,
	}
	logFields = append(logFields, fields...)

	logger.InfoContext(reqCtx, operation+" successful", logFields...)
}

// LogWarning 记录警告日志
func (h *HandlerHelper) LogWarning(ctx *gin.Context, message string, fields ...interface{}) {
	reqCtx := ctx.Request.Context()
	logger.WarnContext(reqCtx, message, fields...)
}

// LogError 记录错误日志
func (h *HandlerHelper) LogError(ctx *gin.Context, message string, fields ...interface{}) {
	reqCtx := ctx.Request.Context()
	logger.ErrorContext(reqCtx, message, fields...)
}

// CheckPermission 检查用户权限（是否为管理员或资源所有者）
func (h *HandlerHelper) CheckPermission(ctx *gin.Context, userID int64, resourceOwnerID int64, operation string) bool {
	reqCtx := ctx.Request.Context()

	// 如果是资源所有者，直接允许
	if userID == resourceOwnerID {
		return true
	}

	// 检查是否是管理员
	isAdmin, exists := ctx.Get("is_admin")
	if exists && cast.ToBool(isAdmin) {
		return true
	}

	// 权限不足
	logger.WarnContext(reqCtx, operation+" permission denied",
		"user_id", userID,
		"resource_owner_id", resourceOwnerID,
		"is_admin", isAdmin)
	response.Error(ctx, errorx.ErrForbidden.New(ctx, errorx.None))
	return false
}

// GetClientInfo 获取客户端信息
func (h *HandlerHelper) GetClientInfo(ctx *gin.Context) (string, string) {
	return ctx.ClientIP(), ctx.Request.UserAgent()
}

// ValidateID 验证ID参数
func (h *HandlerHelper) ValidateID(ctx *gin.Context, idStr string, operation string) (uint, bool) {
	reqCtx := ctx.Request.Context()

	id := cast.ToUint(idStr)
	if id == 0 {
		logger.WarnContext(reqCtx, operation+" invalid ID", "id", idStr)
		response.Error(ctx, errorx.ErrInvalidParams.New(ctx, struct{ Params string }{"ID"}))
		return 0, false
	}

	return id, true
}

// ValidateInt64ID 验证int64类型的ID参数
func (h *HandlerHelper) ValidateInt64ID(ctx *gin.Context, idStr string, operation string) (int64, bool) {
	reqCtx := ctx.Request.Context()

	id := cast.ToInt64(idStr)
	if id == 0 {
		logger.WarnContext(reqCtx, operation+" invalid ID", "id", idStr)
		response.Error(ctx, errorx.ErrInvalidParams.New(ctx, struct{ Params string }{"ID"}))
		return 0, false
	}

	return id, true
}
