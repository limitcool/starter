package controller

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/i18n"
)

// BaseController 基础控制器
type BaseController struct{}

// GetCurrentLanguage 获取当前上下文中的语言
func (b *BaseController) GetCurrentLanguage(c *gin.Context) string {
	return b.GetCurrentLanguageWithContext(c.Request.Context(), c)
}

// GetCurrentLanguageWithContext 使用上下文获取当前上下文中的语言
func (b *BaseController) GetCurrentLanguageWithContext(ctx context.Context, c *gin.Context) string {
	lang, exists := c.Get("lang")
	if !exists {
		return i18n.GetDefaultLanguageContext(ctx)
	}
	return lang.(string)
}
