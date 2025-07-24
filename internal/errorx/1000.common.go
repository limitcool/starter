package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

func init() {
	commonI18n.LoadTranslations()
}

var commonI18n = i18n.NewCatalog("common")

var (
	// 成功
	Success = DefineSimple(commonI18n, 0, "success", http.StatusOK)

	ErrUnknown         = DefineSimple(commonI18n, 5000, "unknown error", http.StatusInternalServerError)                             // 未知错误
	ErrInvalidParams   = Define[struct{ Params string }](commonI18n, 1000, "invalid parameters: {{.Params}}", http.StatusBadRequest) // 请求参数错误
	ErrInternal        = DefineSimple(commonI18n, 1001, "internal error", http.StatusInternalServerError)                            // 服务器内部错误
	ErrUnauthorized    = DefineSimple(commonI18n, 1002, "unauthorized", http.StatusUnauthorized)                                     // 未授权
	ErrForbidden       = DefineSimple(commonI18n, 1003, "forbidden", http.StatusForbidden)                                           // 禁止访问
	ErrNotFound        = DefineSimple(commonI18n, 1004, "resource does not exist", http.StatusNotFound)                              // 资源不存在
	ErrTimeout         = DefineSimple(commonI18n, 1005, "request timeout", http.StatusRequestTimeout)                                // 请求超时
	ErrTooManyRequests = DefineSimple(commonI18n, 1006, "too many requests", http.StatusTooManyRequests)                             // 请求过多
	ErrAccessDenied    = DefineSimple(commonI18n, 1007, "access denied", http.StatusForbidden)                                       // 访问被拒绝
	ErrUserAuthFailed  = DefineSimple(commonI18n, 1008, "user authentication failed", http.StatusUnauthorized)                       // 用户认证失败
	ErrCasbinService   = DefineSimple(commonI18n, 1009, "casbin service error", http.StatusInternalServerError)                      // Casbin服务错误
	ErrFileStorage     = DefineSimple(commonI18n, 1010, "file storage error", http.StatusInternalServerError)                        // 文件存储错误
)
