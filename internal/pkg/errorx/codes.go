package errorx

import "net/http"

// ErrorCode 错误码结构体
type ErrorCode struct {
	Code     int    // 错误码
	Message  string // 错误消息
	HTTPCode int    // HTTP状态码
}

// 预定义错误变量
var (
	// 通用错误
	ErrUnknown         *AppError
	ErrInvalidParams   *AppError
	ErrInternal        *AppError
	ErrUnauthorized    *AppError
	ErrForbidden       *AppError
	ErrNotFound        *AppError
	ErrTimeout         *AppError
	ErrTooManyRequests *AppError
	ErrAccessDenied    *AppError
	ErrUserAuthFailed  *AppError
	ErrCasbinService   *AppError
	ErrFileStorage     *AppError

	// 用户错误
	ErrUserNotFound       *AppError
	ErrInvalidCredentials *AppError
	ErrUserDisabled       *AppError
	ErrUserExists         *AppError
	ErrPasswordExpired    *AppError
	ErrUserPasswordError  *AppError
	ErrUserTokenError     *AppError
	ErrUserNoLogin        *AppError

	// 数据库错误
	ErrDatabase            *AppError
	ErrDatabaseQuery       *AppError
	ErrDatabaseInsert      *AppError
	ErrDatabaseUpdate      *AppError
	ErrDatabaseDelete      *AppError
	ErrDatabaseConnection  *AppError
	ErrDatabaseTransaction *AppError
)

// 错误码常量
var (
	// 基础错误码
	Success = ErrorCode{
		Code:     0,
		Message:  "成功",
		HTTPCode: http.StatusOK,
	}

	// 通用错误码 (1000-1999)
	ErrUnknownCode = ErrorCode{
		Code:     1000,
		Message:  "未知错误",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrInvalidParamsCode = ErrorCode{
		Code:     1001,
		Message:  "无效的参数",
		HTTPCode: http.StatusBadRequest,
	}

	ErrInternalCode = ErrorCode{
		Code:     1002,
		Message:  "内部错误",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrUnauthorizedCode = ErrorCode{
		Code:     1003,
		Message:  "未授权",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrForbiddenCode = ErrorCode{
		Code:     1004,
		Message:  "禁止访问",
		HTTPCode: http.StatusForbidden,
	}

	ErrNotFoundCode = ErrorCode{
		Code:     1005,
		Message:  "资源不存在",
		HTTPCode: http.StatusNotFound,
	}

	ErrTimeoutCode = ErrorCode{
		Code:     1006,
		Message:  "请求超时",
		HTTPCode: http.StatusRequestTimeout,
	}

	ErrTooManyRequestsCode = ErrorCode{
		Code:     1007,
		Message:  "请求过多",
		HTTPCode: http.StatusTooManyRequests,
	}

	ErrAccessDeniedCode = ErrorCode{
		Code:     1008,
		Message:  "访问被拒绝",
		HTTPCode: http.StatusForbidden,
	}

	ErrUserAuthFailedCode = ErrorCode{
		Code:     1009,
		Message:  "用户认证失败",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrCasbinServiceCode = ErrorCode{
		Code:     1010,
		Message:  "Casbin服务错误",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrFileStorageCode = ErrorCode{
		Code:     1011,
		Message:  "文件存储错误",
		HTTPCode: http.StatusInternalServerError,
	}

	// 用户错误码 (2000-2999)
	ErrUserNotFoundCode = ErrorCode{
		Code:     2000,
		Message:  "用户不存在",
		HTTPCode: http.StatusNotFound,
	}
	// 兼容旧代码
	ErrorUserNotFoundCode = ErrorCode{
		Code:     2000,
		Message:  "用户不存在",
		HTTPCode: http.StatusNotFound,
	}

	ErrInvalidCredentialsCode = ErrorCode{
		Code:     2001,
		Message:  "无效的凭证",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrUserDisabledCode = ErrorCode{
		Code:     2002,
		Message:  "用户已禁用",
		HTTPCode: http.StatusForbidden,
	}

	ErrUserExistsCode = ErrorCode{
		Code:     2003,
		Message:  "用户已存在",
		HTTPCode: http.StatusConflict,
	}

	ErrPasswordExpiredCode = ErrorCode{
		Code:     2004,
		Message:  "密码已过期",
		HTTPCode: http.StatusForbidden,
	}

	ErrUserPasswordErrorCode = ErrorCode{
		Code:     2005,
		Message:  "用户密码错误",
		HTTPCode: http.StatusUnauthorized,
	}
	// 兼容旧代码
	ErrorUserPasswordErrorCode = ErrorCode{
		Code:     2005,
		Message:  "用户密码错误",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrUserTokenErrorCode = ErrorCode{
		Code:     2006,
		Message:  "用户令牌错误",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrUserNoLoginCode = ErrorCode{
		Code:     2007,
		Message:  "用户未登录",
		HTTPCode: http.StatusUnauthorized,
	}

	// 数据库错误码 (3000-3999)
	ErrDatabaseCode = ErrorCode{
		Code:     3000,
		Message:  "数据库错误",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrDatabaseQueryCode = ErrorCode{
		Code:     3001,
		Message:  "数据库查询错误",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrDatabaseInsertCode = ErrorCode{
		Code:     3002,
		Message:  "数据库插入错误",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrDatabaseUpdateCode = ErrorCode{
		Code:     3003,
		Message:  "数据库更新错误",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrDatabaseDeleteCode = ErrorCode{
		Code:     3004,
		Message:  "数据库删除错误",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrDatabaseConnectionCode = ErrorCode{
		Code:     3005,
		Message:  "数据库连接错误",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrDatabaseTransactionCode = ErrorCode{
		Code:     3006,
		Message:  "数据库事务错误",
		HTTPCode: http.StatusInternalServerError,
	}

	// 通用错误码别名 (兼容旧代码)
	ErrorNotFoundCode = ErrNotFoundCode
)

// 错误码常量 - 直接使用整数值，方便比较
const (
	// 通用错误码
	ErrUnknownCodeValue         = 1000
	ErrInvalidParamsCodeValue   = 1001
	ErrInternalCodeValue        = 1002
	ErrUnauthorizedCodeValue    = 1003
	ErrForbiddenCodeValue       = 1004
	ErrNotFoundCodeValue        = 1005
	ErrTimeoutCodeValue         = 1006
	ErrTooManyRequestsCodeValue = 1007
	ErrAccessDeniedCodeValue    = 1008
	ErrUserAuthFailedCodeValue  = 1009
	ErrCasbinServiceCodeValue   = 1010
	ErrFileStorageCodeValue     = 1011

	// 用户错误码
	ErrUserNotFoundCodeValue        = 2000
	ErrorUserNotFoundCodeValue      = 2000 // 兼容旧代码
	ErrInvalidCredentialsCodeValue  = 2001
	ErrUserDisabledCodeValue        = 2002
	ErrUserExistsCodeValue          = 2003
	ErrPasswordExpiredCodeValue     = 2004
	ErrUserPasswordErrorCodeValue   = 2005
	ErrorUserPasswordErrorCodeValue = 2005 // 兼容旧代码
	ErrUserTokenErrorCodeValue      = 2006
	ErrUserNoLoginCodeValue         = 2007

	// 数据库错误码
	ErrDatabaseCodeValue            = 3000
	ErrDatabaseQueryCodeValue       = 3001
	ErrDatabaseInsertCodeValue      = 3002
	ErrDatabaseUpdateCodeValue      = 3003
	ErrDatabaseDeleteCodeValue      = 3004
	ErrDatabaseConnectionCodeValue  = 3005
	ErrDatabaseTransactionCodeValue = 3006
)

// 初始化函数，创建预定义错误实例
func init() {
	// 通用错误
	ErrUnknown = NewAppError(ErrUnknownCodeValue, ErrUnknownCode.Message, ErrUnknownCode.HTTPCode)
	ErrInvalidParams = NewAppError(ErrInvalidParamsCodeValue, ErrInvalidParamsCode.Message, ErrInvalidParamsCode.HTTPCode)
	ErrInternal = NewAppError(ErrInternalCodeValue, ErrInternalCode.Message, ErrInternalCode.HTTPCode)
	ErrUnauthorized = NewAppError(ErrUnauthorizedCodeValue, ErrUnauthorizedCode.Message, ErrUnauthorizedCode.HTTPCode)
	ErrForbidden = NewAppError(ErrForbiddenCodeValue, ErrForbiddenCode.Message, ErrForbiddenCode.HTTPCode)
	ErrNotFound = NewAppError(ErrNotFoundCodeValue, ErrNotFoundCode.Message, ErrNotFoundCode.HTTPCode)
	ErrTimeout = NewAppError(ErrTimeoutCodeValue, ErrTimeoutCode.Message, ErrTimeoutCode.HTTPCode)
	ErrTooManyRequests = NewAppError(ErrTooManyRequestsCodeValue, ErrTooManyRequestsCode.Message, ErrTooManyRequestsCode.HTTPCode)
	ErrAccessDenied = NewAppError(ErrAccessDeniedCodeValue, ErrAccessDeniedCode.Message, ErrAccessDeniedCode.HTTPCode)
	ErrUserAuthFailed = NewAppError(ErrUserAuthFailedCodeValue, ErrUserAuthFailedCode.Message, ErrUserAuthFailedCode.HTTPCode)
	ErrCasbinService = NewAppError(ErrCasbinServiceCodeValue, ErrCasbinServiceCode.Message, ErrCasbinServiceCode.HTTPCode)
	ErrFileStorage = NewAppError(ErrFileStorageCodeValue, ErrFileStorageCode.Message, ErrFileStorageCode.HTTPCode)

	// 用户错误
	ErrUserNotFound = NewAppError(ErrUserNotFoundCodeValue, ErrUserNotFoundCode.Message, ErrUserNotFoundCode.HTTPCode)
	ErrInvalidCredentials = NewAppError(ErrInvalidCredentialsCodeValue, ErrInvalidCredentialsCode.Message, ErrInvalidCredentialsCode.HTTPCode)
	ErrUserDisabled = NewAppError(ErrUserDisabledCodeValue, ErrUserDisabledCode.Message, ErrUserDisabledCode.HTTPCode)
	ErrUserExists = NewAppError(ErrUserExistsCodeValue, ErrUserExistsCode.Message, ErrUserExistsCode.HTTPCode)
	ErrPasswordExpired = NewAppError(ErrPasswordExpiredCodeValue, ErrPasswordExpiredCode.Message, ErrPasswordExpiredCode.HTTPCode)
	ErrUserPasswordError = NewAppError(ErrUserPasswordErrorCodeValue, ErrUserPasswordErrorCode.Message, ErrUserPasswordErrorCode.HTTPCode)
	ErrUserTokenError = NewAppError(ErrUserTokenErrorCodeValue, ErrUserTokenErrorCode.Message, ErrUserTokenErrorCode.HTTPCode)
	ErrUserNoLogin = NewAppError(ErrUserNoLoginCodeValue, ErrUserNoLoginCode.Message, ErrUserNoLoginCode.HTTPCode)

	// 数据库错误
	ErrDatabase = NewAppError(ErrDatabaseCodeValue, ErrDatabaseCode.Message, ErrDatabaseCode.HTTPCode)
	ErrDatabaseQuery = NewAppError(ErrDatabaseQueryCodeValue, ErrDatabaseQueryCode.Message, ErrDatabaseQueryCode.HTTPCode)
	ErrDatabaseInsert = NewAppError(ErrDatabaseInsertCodeValue, ErrDatabaseInsertCode.Message, ErrDatabaseInsertCode.HTTPCode)
	ErrDatabaseUpdate = NewAppError(ErrDatabaseUpdateCodeValue, ErrDatabaseUpdateCode.Message, ErrDatabaseUpdateCode.HTTPCode)
	ErrDatabaseDelete = NewAppError(ErrDatabaseDeleteCodeValue, ErrDatabaseDeleteCode.Message, ErrDatabaseDeleteCode.HTTPCode)
	ErrDatabaseConnection = NewAppError(ErrDatabaseConnectionCodeValue, ErrDatabaseConnectionCode.Message, ErrDatabaseConnectionCode.HTTPCode)
	ErrDatabaseTransaction = NewAppError(ErrDatabaseTransactionCodeValue, ErrDatabaseTransactionCode.Message, ErrDatabaseTransactionCode.HTTPCode)
}

// 权限错误码 (4000-4999)
var (
	ErrPermissionDeniedCode = ErrorCode{
		Code:     4000,
		Message:  "权限不足",
		HTTPCode: http.StatusForbidden,
	}

	ErrRoleNotFoundCode = ErrorCode{
		Code:     4001,
		Message:  "角色不存在",
		HTTPCode: http.StatusNotFound,
	}

	ErrRoleExistsCode = ErrorCode{
		Code:     4002,
		Message:  "角色已存在",
		HTTPCode: http.StatusConflict,
	}

	ErrMenuNotFoundCode = ErrorCode{
		Code:     4003,
		Message:  "菜单不存在",
		HTTPCode: http.StatusNotFound,
	}

	ErrAPINotFoundCode = ErrorCode{
		Code:     4004,
		Message:  "API不存在",
		HTTPCode: http.StatusNotFound,
	}
)

// 文件错误码 (5000-5999)
var (
	ErrFileNotFoundCode = ErrorCode{
		Code:     5000,
		Message:  "文件不存在",
		HTTPCode: http.StatusNotFound,
	}

	ErrFileUploadCode = ErrorCode{
		Code:     5001,
		Message:  "文件上传失败",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrFileDownloadCode = ErrorCode{
		Code:     5002,
		Message:  "文件下载失败",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrFileDeleteCode = ErrorCode{
		Code:     5003,
		Message:  "文件删除失败",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrFileSizeExceededCode = ErrorCode{
		Code:     5004,
		Message:  "文件大小超出限制",
		HTTPCode: http.StatusBadRequest,
	}

	ErrFileTypeNotAllowedCode = ErrorCode{
		Code:     5005,
		Message:  "文件类型不允许",
		HTTPCode: http.StatusBadRequest,
	}
)

// 缓存错误码 (6000-6999)
var (
	ErrCacheCode = ErrorCode{
		Code:     6000,
		Message:  "缓存错误",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrCacheKeyNotFoundCode = ErrorCode{
		Code:     6001,
		Message:  "缓存键不存在",
		HTTPCode: http.StatusNotFound,
	}

	ErrCacheSetCode = ErrorCode{
		Code:     6002,
		Message:  "缓存设置失败",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrCacheDeleteCode = ErrorCode{
		Code:     6003,
		Message:  "缓存删除失败",
		HTTPCode: http.StatusInternalServerError,
	}
)

// 错误码映射表
var errorCodeMap map[int]ErrorCode

// 创建预定义错误实例
func init() {
	// 初始化错误码映射表
	errorCodeMap = make(map[int]ErrorCode)

	// 添加基础错误码
	errorCodeMap[Success.Code] = Success

	// 添加通用错误码
	errorCodeMap[ErrUnknownCode.Code] = ErrUnknownCode
	errorCodeMap[ErrInvalidParamsCode.Code] = ErrInvalidParamsCode
	errorCodeMap[ErrInternalCode.Code] = ErrInternalCode
	errorCodeMap[ErrUnauthorizedCode.Code] = ErrUnauthorizedCode
	errorCodeMap[ErrForbiddenCode.Code] = ErrForbiddenCode
	errorCodeMap[ErrNotFoundCode.Code] = ErrNotFoundCode
	errorCodeMap[ErrTimeoutCode.Code] = ErrTimeoutCode
	errorCodeMap[ErrTooManyRequestsCode.Code] = ErrTooManyRequestsCode
	errorCodeMap[ErrAccessDeniedCode.Code] = ErrAccessDeniedCode
	errorCodeMap[ErrUserAuthFailedCode.Code] = ErrUserAuthFailedCode
	errorCodeMap[ErrCasbinServiceCode.Code] = ErrCasbinServiceCode
	errorCodeMap[ErrFileStorageCode.Code] = ErrFileStorageCode

	// 添加用户错误码
	errorCodeMap[ErrUserNotFoundCode.Code] = ErrUserNotFoundCode
	errorCodeMap[ErrInvalidCredentialsCode.Code] = ErrInvalidCredentialsCode
	errorCodeMap[ErrUserDisabledCode.Code] = ErrUserDisabledCode
	errorCodeMap[ErrUserExistsCode.Code] = ErrUserExistsCode
	errorCodeMap[ErrPasswordExpiredCode.Code] = ErrPasswordExpiredCode
	errorCodeMap[ErrUserPasswordErrorCode.Code] = ErrUserPasswordErrorCode
	errorCodeMap[ErrUserTokenErrorCode.Code] = ErrUserTokenErrorCode
	errorCodeMap[ErrUserNoLoginCode.Code] = ErrUserNoLoginCode

	// 添加数据库错误码
	errorCodeMap[ErrDatabaseCode.Code] = ErrDatabaseCode
	errorCodeMap[ErrDatabaseQueryCode.Code] = ErrDatabaseQueryCode
	errorCodeMap[ErrDatabaseInsertCode.Code] = ErrDatabaseInsertCode
	errorCodeMap[ErrDatabaseUpdateCode.Code] = ErrDatabaseUpdateCode
	errorCodeMap[ErrDatabaseDeleteCode.Code] = ErrDatabaseDeleteCode
	errorCodeMap[ErrDatabaseConnectionCode.Code] = ErrDatabaseConnectionCode
	errorCodeMap[ErrDatabaseTransactionCode.Code] = ErrDatabaseTransactionCode

	// 添加权限错误码
	errorCodeMap[ErrPermissionDeniedCode.Code] = ErrPermissionDeniedCode
	errorCodeMap[ErrRoleNotFoundCode.Code] = ErrRoleNotFoundCode
	errorCodeMap[ErrRoleExistsCode.Code] = ErrRoleExistsCode
	errorCodeMap[ErrMenuNotFoundCode.Code] = ErrMenuNotFoundCode
	errorCodeMap[ErrAPINotFoundCode.Code] = ErrAPINotFoundCode

	// 添加文件错误码
	errorCodeMap[ErrFileNotFoundCode.Code] = ErrFileNotFoundCode
	errorCodeMap[ErrFileUploadCode.Code] = ErrFileUploadCode
	errorCodeMap[ErrFileDownloadCode.Code] = ErrFileDownloadCode
	errorCodeMap[ErrFileDeleteCode.Code] = ErrFileDeleteCode
	errorCodeMap[ErrFileSizeExceededCode.Code] = ErrFileSizeExceededCode
	errorCodeMap[ErrFileTypeNotAllowedCode.Code] = ErrFileTypeNotAllowedCode

	// 添加缓存错误码
	errorCodeMap[ErrCacheCode.Code] = ErrCacheCode
	errorCodeMap[ErrCacheKeyNotFoundCode.Code] = ErrCacheKeyNotFoundCode
	errorCodeMap[ErrCacheSetCode.Code] = ErrCacheSetCode
	errorCodeMap[ErrCacheDeleteCode.Code] = ErrCacheDeleteCode

	// 创建预定义错误实例
	// 通用错误
	ErrUnknown = NewAppError(ErrUnknownCode.Code, ErrUnknownCode.Message, ErrUnknownCode.HTTPCode)
	ErrInvalidParams = NewAppError(ErrInvalidParamsCode.Code, ErrInvalidParamsCode.Message, ErrInvalidParamsCode.HTTPCode)
	ErrInternal = NewAppError(ErrInternalCode.Code, ErrInternalCode.Message, ErrInternalCode.HTTPCode)
	ErrUnauthorized = NewAppError(ErrUnauthorizedCode.Code, ErrUnauthorizedCode.Message, ErrUnauthorizedCode.HTTPCode)
	ErrForbidden = NewAppError(ErrForbiddenCode.Code, ErrForbiddenCode.Message, ErrForbiddenCode.HTTPCode)
	ErrNotFound = NewAppError(ErrNotFoundCode.Code, ErrNotFoundCode.Message, ErrNotFoundCode.HTTPCode)
	ErrTimeout = NewAppError(ErrTimeoutCode.Code, ErrTimeoutCode.Message, ErrTimeoutCode.HTTPCode)
	ErrTooManyRequests = NewAppError(ErrTooManyRequestsCode.Code, ErrTooManyRequestsCode.Message, ErrTooManyRequestsCode.HTTPCode)
	ErrAccessDenied = NewAppError(ErrAccessDeniedCode.Code, ErrAccessDeniedCode.Message, ErrAccessDeniedCode.HTTPCode)
	ErrUserAuthFailed = NewAppError(ErrUserAuthFailedCode.Code, ErrUserAuthFailedCode.Message, ErrUserAuthFailedCode.HTTPCode)
	ErrCasbinService = NewAppError(ErrCasbinServiceCode.Code, ErrCasbinServiceCode.Message, ErrCasbinServiceCode.HTTPCode)
	ErrFileStorage = NewAppError(ErrFileStorageCode.Code, ErrFileStorageCode.Message, ErrFileStorageCode.HTTPCode)

	// 用户错误
	ErrUserNotFound = NewAppError(ErrUserNotFoundCode.Code, ErrUserNotFoundCode.Message, ErrUserNotFoundCode.HTTPCode)
	ErrInvalidCredentials = NewAppError(ErrInvalidCredentialsCode.Code, ErrInvalidCredentialsCode.Message, ErrInvalidCredentialsCode.HTTPCode)
	ErrUserDisabled = NewAppError(ErrUserDisabledCode.Code, ErrUserDisabledCode.Message, ErrUserDisabledCode.HTTPCode)
	ErrUserExists = NewAppError(ErrUserExistsCode.Code, ErrUserExistsCode.Message, ErrUserExistsCode.HTTPCode)
	ErrPasswordExpired = NewAppError(ErrPasswordExpiredCode.Code, ErrPasswordExpiredCode.Message, ErrPasswordExpiredCode.HTTPCode)
	ErrUserPasswordError = NewAppError(ErrUserPasswordErrorCode.Code, ErrUserPasswordErrorCode.Message, ErrUserPasswordErrorCode.HTTPCode)
	ErrUserTokenError = NewAppError(ErrUserTokenErrorCode.Code, ErrUserTokenErrorCode.Message, ErrUserTokenErrorCode.HTTPCode)
	ErrUserNoLogin = NewAppError(ErrUserNoLoginCode.Code, ErrUserNoLoginCode.Message, ErrUserNoLoginCode.HTTPCode)

	// 数据库错误
	ErrDatabase = NewAppError(ErrDatabaseCode.Code, ErrDatabaseCode.Message, ErrDatabaseCode.HTTPCode)
	ErrDatabaseQuery = NewAppError(ErrDatabaseQueryCode.Code, ErrDatabaseQueryCode.Message, ErrDatabaseQueryCode.HTTPCode)
	ErrDatabaseInsert = NewAppError(ErrDatabaseInsertCode.Code, ErrDatabaseInsertCode.Message, ErrDatabaseInsertCode.HTTPCode)
	ErrDatabaseUpdate = NewAppError(ErrDatabaseUpdateCode.Code, ErrDatabaseUpdateCode.Message, ErrDatabaseUpdateCode.HTTPCode)
	ErrDatabaseDelete = NewAppError(ErrDatabaseDeleteCode.Code, ErrDatabaseDeleteCode.Message, ErrDatabaseDeleteCode.HTTPCode)
	ErrDatabaseConnection = NewAppError(ErrDatabaseConnectionCode.Code, ErrDatabaseConnectionCode.Message, ErrDatabaseConnectionCode.HTTPCode)
	ErrDatabaseTransaction = NewAppError(ErrDatabaseTransactionCode.Code, ErrDatabaseTransactionCode.Message, ErrDatabaseTransactionCode.HTTPCode)
}

// FindErrorCode 根据错误码获取错误码结构体
func FindErrorCode(code int) ErrorCode {
	if errCode, ok := errorCodeMap[code]; ok {
		return errCode
	}
	return ErrUnknownCode
}

// NewErrorFromCode 根据错误码创建新的错误实例
func NewErrorFromCode(code int) *AppError {
	errCode := FindErrorCode(code)
	return NewAppError(errCode.Code, errCode.Message, errCode.HTTPCode)
}
