package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/i18n"
)

// 语言标识符的来源
const (
	headerAcceptLanguage = "Accept-Language"
	queryLang            = "lang"
	cookieLang           = "lang"
)

// I18n 国际化中间件
func I18n(config *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取配置
		if !config.I18n.Enabled {
			c.Next()
			return
		}

		// 获取语言标识
		// 查找顺序：1. URL参数 2. Cookie 3. Accept-Language头 4. 默认语言
		lang := c.Query(queryLang)
		if lang == "" {
			// 尝试从cookie中获取
			langCookie, _ := c.Cookie(cookieLang)
			if langCookie != "" {
				lang = langCookie
			} else {
				// 尝试从Accept-Language获取
				acceptLanguage := c.GetHeader(headerAcceptLanguage)
				if acceptLanguage != "" {
					// 解析Accept-Language
					preferredLang := parseAcceptLanguage(acceptLanguage, config.I18n.SupportLanguages)
					if preferredLang != "" {
						lang = preferredLang
					}
				}
			}
		}

		// 如果没有获取到语言或语言不受支持，使用默认语言
		if lang == "" || !isSupported(lang, config.I18n.SupportLanguages) {
			lang = config.I18n.DefaultLanguage
		}

		// 将语言设置到上下文
		c.Set("lang", lang)

		// 如果请求是通过URL参数设置语言的，将其保存到cookie
		if c.Query(queryLang) != "" {
			c.SetCookie(cookieLang, lang, 60*60*24*30, "/", "", false, false)
		}

		c.Next()
	}
}

// 解析Accept-Language头并返回最佳匹配的语言
func parseAcceptLanguage(acceptLanguage string, supportedLangs []string) string {
	// 简单解析Accept-Language
	languages := strings.Split(acceptLanguage, ",")
	if len(languages) == 0 {
		return ""
	}

	// 提取语言代码（去除权重部分）
	for _, lang := range languages {
		langCode := strings.Split(strings.TrimSpace(lang), ";")[0]
		for _, supported := range supportedLangs {
			// 检查是否匹配支持的语言
			if strings.EqualFold(langCode, supported) ||
				strings.HasPrefix(strings.ToLower(langCode), strings.ToLower(supported[:2])) {
				return supported
			}
		}
	}

	return ""
}

// isSupported 检查语言是否在支持列表中
func isSupported(lang string, supportedLangs []string) bool {
	for _, supported := range supportedLangs {
		if lang == supported {
			return true
		}
	}
	return false
}

// TranslateError 翻译错误消息的辅助函数
func TranslateError(c *gin.Context, messageID string) string {
	return TranslateErrorWithContext(c.Request.Context(), c, messageID)
}

// TranslateErrorWithContext 使用上下文翻译错误消息的辅助函数
func TranslateErrorWithContext(ctx context.Context, c *gin.Context, messageID string) string {
	lang, exists := c.Get("lang")
	if !exists {
		// 如果上下文中没有语言信息，使用默认语言
		return i18n.TContext(ctx, messageID, i18n.GetDefaultLanguageContext(ctx))
	}
	return i18n.TContext(ctx, messageID, lang.(string))
}
