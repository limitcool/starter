package errorx

import "net/http"

// 基础错误码
const (
	SuccessCode      = 0      // 成功
	ErrorUnknownCode = 500000 // 未知错误
)

// 通用错误码 (1000-1999)
const (
	ErrorInvalidParamsCode   = 1000 // 无效的参数
	ErrorInternalCode        = 1001 // 内部错误
	ErrorUnauthorizedCode    = 1002 // 未授权
	ErrorForbiddenCode       = 1003 // 禁止访问
	ErrorNotFoundCode        = 1004 // 资源不存在
	ErrorTimeoutCode         = 1005 // 请求超时
	ErrorTooManyRequestsCode = 1006 // 请求过多
	ErrorAccessDeniedCode    = 1007 // 访问被拒绝
	ErrorUserAuthFailedCode  = 1008 // 用户认证失败
	ErrorCasbinServiceCode   = 1009 // Casbin服务错误
	ErrorFileStorageCode     = 1010 // 文件存储错误
)

// 用户错误码 (2000-2999)
const (
	ErrorUserNotFoundCode       = 2000 // 用户不存在
	ErrorInvalidCredentialsCode = 2001 // 无效的凭证
	ErrorUserDisabledCode       = 2002 // 用户已禁用
	ErrorUserExistsCode         = 2003 // 用户已存在
	ErrorPasswordExpiredCode    = 2004 // 密码已过期
	ErrorUserPasswordErrorCode  = 2005 // 用户密码错误
	ErrorUserTokenErrorCode     = 2006 // 用户令牌错误
	ErrorUserNoLoginCode        = 2007 // 用户未登录
)

// 数据库错误码 (3000-3999)
const (
	ErrorDatabaseCode            = 3000 // 数据库错误
	ErrorDatabaseQueryCode       = 3001 // 数据库查询错误
	ErrorDatabaseInsertCode      = 3002 // 数据库插入错误
	ErrorDatabaseUpdateCode      = 3003 // 数据库更新错误
	ErrorDatabaseDeleteCode      = 3004 // 数据库删除错误
	ErrorDatabaseConnectionCode  = 3005 // 数据库连接错误
	ErrorDatabaseTransactionCode = 3006 // 数据库事务错误
)

// 预定义错误实例
var (
	// 成功
	Success = NewAppError(SuccessCode, "成功", http.StatusOK)

	// 通用错误
	ErrUnknown         = NewAppError(ErrorUnknownCode, "未知错误", http.StatusInternalServerError)
	ErrInvalidParams   = NewAppError(ErrorInvalidParamsCode, "无效的参数", http.StatusBadRequest)
	ErrInternal        = NewAppError(ErrorInternalCode, "内部错误", http.StatusInternalServerError)
	ErrUnauthorized    = NewAppError(ErrorUnauthorizedCode, "未授权", http.StatusUnauthorized)
	ErrForbidden       = NewAppError(ErrorForbiddenCode, "禁止访问", http.StatusForbidden)
	ErrNotFound        = NewAppError(ErrorNotFoundCode, "资源不存在", http.StatusNotFound)
	ErrTimeout         = NewAppError(ErrorTimeoutCode, "请求超时", http.StatusRequestTimeout)
	ErrTooManyRequests = NewAppError(ErrorTooManyRequestsCode, "请求过多", http.StatusTooManyRequests)
	ErrAccessDenied    = NewAppError(ErrorAccessDeniedCode, "访问被拒绝", http.StatusForbidden)
	ErrUserAuthFailed  = NewAppError(ErrorUserAuthFailedCode, "用户认证失败", http.StatusUnauthorized)
	ErrCasbinService   = NewAppError(ErrorCasbinServiceCode, "Casbin服务错误", http.StatusInternalServerError)
	ErrFileStorage     = NewAppError(ErrorFileStorageCode, "文件存储错误", http.StatusInternalServerError)

	// 用户错误
	ErrUserNotFound       = NewAppError(ErrorUserNotFoundCode, "用户不存在", http.StatusNotFound)
	ErrInvalidCredentials = NewAppError(ErrorInvalidCredentialsCode, "无效的凭证", http.StatusUnauthorized)
	ErrUserDisabled       = NewAppError(ErrorUserDisabledCode, "用户已禁用", http.StatusForbidden)
	ErrUserExists         = NewAppError(ErrorUserExistsCode, "用户已存在", http.StatusConflict)
	ErrPasswordExpired    = NewAppError(ErrorPasswordExpiredCode, "密码已过期", http.StatusForbidden)
	ErrUserPasswordError  = NewAppError(ErrorUserPasswordErrorCode, "用户密码错误", http.StatusUnauthorized)
	ErrUserTokenError     = NewAppError(ErrorUserTokenErrorCode, "用户令牌错误", http.StatusUnauthorized)
	ErrUserNoLogin        = NewAppError(ErrorUserNoLoginCode, "用户未登录", http.StatusUnauthorized)

	// 数据库错误
	ErrDatabase            = NewAppError(ErrorDatabaseCode, "数据库错误", http.StatusInternalServerError)
	ErrDatabaseQuery       = NewAppError(ErrorDatabaseQueryCode, "数据库查询错误", http.StatusInternalServerError)
	ErrDatabaseInsert      = NewAppError(ErrorDatabaseInsertCode, "数据库插入错误", http.StatusInternalServerError)
	ErrDatabaseUpdate      = NewAppError(ErrorDatabaseUpdateCode, "数据库更新错误", http.StatusInternalServerError)
	ErrDatabaseDelete      = NewAppError(ErrorDatabaseDeleteCode, "数据库删除错误", http.StatusInternalServerError)
	ErrDatabaseConnection  = NewAppError(ErrorDatabaseConnectionCode, "数据库连接错误", http.StatusInternalServerError)
	ErrDatabaseTransaction = NewAppError(ErrorDatabaseTransactionCode, "数据库事务错误", http.StatusInternalServerError)
)
