package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/i18n"
)

// BaseController 基础控制器
type BaseController struct{}

// GetCurrentLanguage 获取当前上下文中的语言
func (b *BaseController) GetCurrentLanguage(c *gin.Context) string {
	lang, exists := c.Get("lang")
	if !exists {
		return i18n.GetDefaultLanguage()
	}
	return lang.(string)
}
