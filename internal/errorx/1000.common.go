package errorx

import (
	"context"
	"net/http"

	"github.com/epkgs/i18n"
)

func init() {
	commonI18n.LoadTranslations()

	ErrFileTest.New(context.Background())
}

var commonI18n = i18n.NewCatalog("common")

var (
	// 成功
	Success = defineErrSimple(commonI18n, 0, "success", http.StatusOK)

	ErrUnknown         = defineErrSimple(commonI18n, 500000, "unknown error", http.StatusInternalServerError)                           // 未知错误
	ErrInvalidParams   = defineErr[struct{ Params string }](commonI18n, 1000, "invalid parameters: {{.Params}}", http.StatusBadRequest) // 请求参数错误
	ErrInternal        = defineErrSimple(commonI18n, 1001, "internal error", http.StatusInternalServerError)                            // 服务器内部错误
	ErrUnauthorized    = defineErrSimple(commonI18n, 1002, "unauthorized", http.StatusUnauthorized)                                     // 未授权
	ErrForbidden       = defineErrSimple(commonI18n, 1003, "forbidden", http.StatusForbidden)                                           // 禁止访问
	ErrNotFound        = defineErrSimple(commonI18n, 1004, "resource does not exist", http.StatusNotFound)                              // 资源不存在
	ErrTimeout         = defineErrSimple(commonI18n, 1005, "request timeout", http.StatusRequestTimeout)                                // 请求超时
	ErrTooManyRequests = defineErrSimple(commonI18n, 1006, "too many requests", http.StatusTooManyRequests)                             // 请求过多
	ErrAccessDenied    = defineErrSimple(commonI18n, 1007, "access denied", http.StatusForbidden)                                       // 访问被拒绝
	ErrUserAuthFailed  = defineErrSimple(commonI18n, 1008, "user authentication failed", http.StatusUnauthorized)                       // 用户认证失败
	ErrCasbinService   = defineErrSimple(commonI18n, 1009, "casbin service error", http.StatusInternalServerError)                      // Casbin服务错误
	ErrFileStorage     = defineErrSimple(commonI18n, 1010, "file storage error", http.StatusInternalServerError)                        // 文件存储错误

	ErrFileTest = defineErrSimple(commonI18n, 1011, "file test error", http.StatusInternalServerError)
)
