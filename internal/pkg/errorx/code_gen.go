// 本文件由工具自动生成，请勿手动修改
// 生成命令: go run tools/errorgen/main.go

package errorx

import "net/http"

// 错误码基础值定义
const (
	// 错误码基础值 (按模块划分)
	CommonErrorBase  = 10000 // 通用错误码 (10000-19999)
	DatabaseErrorBase = 20000 // 数据库错误码 (20000-29999)
	UserErrorBase = 30000 // 用户相关错误码 (30000-39999)
	AuthErrorBase = 40000 // 权限相关错误码 (40000-49999)
	CacheErrorBase = 50000 // 缓存相关错误码 (50000-59999)
	FileErrorBase = 60000 // 文件相关错误码 (60000-69999)
)

// 错误码定义
const (
	// 基础错误码
	SuccessCode = 0 // 成功
	// 通用错误码
	InvalidParamsCode = CommonErrorBase + 0 // 请求参数错误
	ErrorUnknownCode = CommonErrorBase + 1 // 服务器开小差啦，稍后再来试一试
	ErrorNotExistCertCode = CommonErrorBase + 2 // 不存在的认证类型
	ErrorNotFoundCode = CommonErrorBase + 3 // 资源不存在
	ErrorDatabaseCode = CommonErrorBase + 4 // 数据库操作失败
	ErrorInternalCode = CommonErrorBase + 5 // 服务器内部错误
	ErrorCode = CommonErrorBase + 6 // 错误
	ErrorParamCode = CommonErrorBase + 7 // 参数错误
	// 数据库错误码
	DatabaseInsertErrorCode = DatabaseErrorBase + 0 // 数据库插入失败
	DatabaseDeleteErrorCode = DatabaseErrorBase + 1 // 数据库删除失败
	DatabaseQueryErrorCode = DatabaseErrorBase + 2 // 数据库查询失败
	// 用户相关错误码
	UserNoLoginCode = UserErrorBase + 0 // 用户未登录
	UserNotFoundCode = UserErrorBase + 1 // 用户不存在
	UserPasswordErrorCode = UserErrorBase + 2 // 密码错误
	UserNotVerifyCode = UserErrorBase + 3 // 用户未验证
	UserLockedCode = UserErrorBase + 4 // 用户已锁定
	UserDisabledCode = UserErrorBase + 5 // 用户已禁用
	UserExpiredCode = UserErrorBase + 6 // 用户已过期
	UserAlreadyExistsCode = UserErrorBase + 7 // 用户已存在
	UserNameOrPasswordErrorCode = UserErrorBase + 8 // 用户名或密码错误
	UserAuthFailedCode = UserErrorBase + 9 // 认证失败
	UserNoPermissionCode = UserErrorBase + 10 // 没有权限
	UserPasswordErrCode = UserErrorBase + 11 // 密码错误
	UserNotExistCode = UserErrorBase + 12 // 用户不存在
	UserTokenErrorCode = UserErrorBase + 13 // 登录凭证无效
	UserTokenExpiredCode = UserErrorBase + 14 // 登录已过期，请重新登录
	// 权限相关错误码
	AccessDeniedCode = AuthErrorBase + 0 // 访问被拒绝
	CasbinServiceCode = AuthErrorBase + 1 // Casbin服务错误
	// 缓存相关错误码
	ErrorCacheCode = CacheErrorBase + 0 // 缓存操作失败
	ErrorCacheTimeoutCode = CacheErrorBase + 1 // 缓存操作超时
	ErrorCacheKeyCode = CacheErrorBase + 2 // 缓存键不存在
	ErrorCacheValueCode = CacheErrorBase + 3 // 缓存值错误
	// 文件相关错误码
	FileStroageCode = FileErrorBase + 0 // 文件存储失败
)

// 预定义错误实例
var (
	// 基础错误
	ErrSuccess = NewAppError(SuccessCode, "成功", http.StatusOK)
	// 通用错误码实例
	ErrInvalidParams = NewAppError(InvalidParamsCode, "请求参数错误", http.StatusBadRequest)
	ErrUnknown = NewAppError(ErrorUnknownCode, "服务器开小差啦，稍后再来试一试", http.StatusInternalServerError)
	ErrNotExistCert = NewAppError(ErrorNotExistCertCode, "不存在的认证类型", http.StatusBadRequest)
	ErrNotFound = NewAppError(ErrorNotFoundCode, "资源不存在", http.StatusNotFound)
	ErrDatabase = NewAppError(ErrorDatabaseCode, "数据库操作失败", http.StatusInternalServerError)
	ErrInternal = NewAppError(ErrorInternalCode, "服务器内部错误", http.StatusInternalServerError)
	ErrError = NewAppError(ErrorCode, "错误", http.StatusInternalServerError)
	ErrParam = NewAppError(ErrorParamCode, "参数错误", http.StatusBadRequest)
	// 数据库错误码实例
	ErrDatabaseInsertError = NewAppError(DatabaseInsertErrorCode, "数据库插入失败", http.StatusInternalServerError)
	ErrDatabaseDeleteError = NewAppError(DatabaseDeleteErrorCode, "数据库删除失败", http.StatusInternalServerError)
	ErrDatabaseQueryError = NewAppError(DatabaseQueryErrorCode, "数据库查询失败", http.StatusInternalServerError)
	// 用户相关错误码实例
	ErrUserNoLogin = NewAppError(UserNoLoginCode, "用户未登录", http.StatusUnauthorized)
	ErrUserNotFound = NewAppError(UserNotFoundCode, "用户不存在", http.StatusNotFound)
	ErrUserPasswordError = NewAppError(UserPasswordErrorCode, "密码错误", http.StatusUnauthorized)
	ErrUserNotVerify = NewAppError(UserNotVerifyCode, "用户未验证", http.StatusUnauthorized)
	ErrUserLocked = NewAppError(UserLockedCode, "用户已锁定", http.StatusUnauthorized)
	ErrUserDisabled = NewAppError(UserDisabledCode, "用户已禁用", http.StatusUnauthorized)
	ErrUserExpired = NewAppError(UserExpiredCode, "用户已过期", http.StatusUnauthorized)
	ErrUserAlreadyExists = NewAppError(UserAlreadyExistsCode, "用户已存在", http.StatusUnauthorized)
	ErrUserNameOrPasswordError = NewAppError(UserNameOrPasswordErrorCode, "用户名或密码错误", http.StatusUnauthorized)
	ErrUserAuthFailed = NewAppError(UserAuthFailedCode, "认证失败", http.StatusUnauthorized)
	ErrUserNoPermission = NewAppError(UserNoPermissionCode, "没有权限", http.StatusUnauthorized)
	ErrUserPasswordErr = NewAppError(UserPasswordErrCode, "密码错误", http.StatusUnauthorized)
	ErrUserNotExist = NewAppError(UserNotExistCode, "用户不存在", http.StatusUnauthorized)
	ErrUserTokenError = NewAppError(UserTokenErrorCode, "登录凭证无效", http.StatusUnauthorized)
	ErrUserTokenExpired = NewAppError(UserTokenExpiredCode, "登录已过期，请重新登录", http.StatusUnauthorized)
	// 权限相关错误码实例
	ErrAccessDenied = NewAppError(AccessDeniedCode, "访问被拒绝", http.StatusForbidden)
	ErrCasbinService = NewAppError(CasbinServiceCode, "Casbin服务错误", http.StatusInternalServerError)
	// 缓存相关错误码实例
	ErrCache = NewAppError(ErrorCacheCode, "缓存操作失败", http.StatusInternalServerError)
	ErrCacheTimeout = NewAppError(ErrorCacheTimeoutCode, "缓存操作超时", http.StatusInternalServerError)
	ErrCacheKey = NewAppError(ErrorCacheKeyCode, "缓存键不存在", http.StatusInternalServerError)
	ErrCacheValue = NewAppError(ErrorCacheValueCode, "缓存值错误", http.StatusInternalServerError)
	// 文件相关错误码实例
	ErrFileStroage = NewAppError(FileStroageCode, "文件存储失败", http.StatusInternalServerError)
)

// GetError 根据错误码获取预定义错误
func GetError(code int) *AppError {
	switch code {
	case SuccessCode:
		return ErrSuccess
	case InvalidParamsCode:
		return ErrInvalidParams
	case ErrorUnknownCode:
		return ErrUnknown
	case ErrorNotExistCertCode:
		return ErrNotExistCert
	case ErrorNotFoundCode:
		return ErrNotFound
	case ErrorDatabaseCode:
		return ErrDatabase
	case ErrorInternalCode:
		return ErrInternal
	case ErrorCode:
		return ErrError
	case ErrorParamCode:
		return ErrParam
	case DatabaseInsertErrorCode:
		return ErrDatabaseInsertError
	case DatabaseDeleteErrorCode:
		return ErrDatabaseDeleteError
	case DatabaseQueryErrorCode:
		return ErrDatabaseQueryError
	case UserNoLoginCode:
		return ErrUserNoLogin
	case UserNotFoundCode:
		return ErrUserNotFound
	case UserPasswordErrorCode:
		return ErrUserPasswordError
	case UserNotVerifyCode:
		return ErrUserNotVerify
	case UserLockedCode:
		return ErrUserLocked
	case UserDisabledCode:
		return ErrUserDisabled
	case UserExpiredCode:
		return ErrUserExpired
	case UserAlreadyExistsCode:
		return ErrUserAlreadyExists
	case UserNameOrPasswordErrorCode:
		return ErrUserNameOrPasswordError
	case UserAuthFailedCode:
		return ErrUserAuthFailed
	case UserNoPermissionCode:
		return ErrUserNoPermission
	case UserPasswordErrCode:
		return ErrUserPasswordErr
	case UserNotExistCode:
		return ErrUserNotExist
	case UserTokenErrorCode:
		return ErrUserTokenError
	case UserTokenExpiredCode:
		return ErrUserTokenExpired
	case AccessDeniedCode:
		return ErrAccessDenied
	case CasbinServiceCode:
		return ErrCasbinService
	case ErrorCacheCode:
		return ErrCache
	case ErrorCacheTimeoutCode:
		return ErrCacheTimeout
	case ErrorCacheKeyCode:
		return ErrCacheKey
	case ErrorCacheValueCode:
		return ErrCacheValue
	case FileStroageCode:
		return ErrFileStroage
	default:
		return ErrUnknown // 默认返回未知错误
	}
}