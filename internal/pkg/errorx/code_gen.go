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
	InvalidParamsCode = CommonErrorBase + 0 // Invalid request parameters
	ErrorUnknownCode = CommonErrorBase + 1 // Server is busy, please try again later
	ErrorNotExistCertCode = CommonErrorBase + 2 // Authentication type does not exist
	ErrorNotFoundCode = CommonErrorBase + 3 // Resource not found
	ErrorDatabaseCode = CommonErrorBase + 4 // Database operation failed
	ErrorInternalCode = CommonErrorBase + 5 // Internal server error
	ErrorCode = CommonErrorBase + 6 // Error
	ErrorParamCode = CommonErrorBase + 7 // Parameter error
	// 数据库错误码
	DatabaseInsertErrorCode = DatabaseErrorBase + 0 // Database insert failed
	DatabaseDeleteErrorCode = DatabaseErrorBase + 1 // Database delete failed
	DatabaseQueryErrorCode = DatabaseErrorBase + 2 // Database query failed
	// 用户相关错误码
	UserNoLoginCode = UserErrorBase + 0 // User not logged in
	UserNotFoundCode = UserErrorBase + 1 // User not found
	UserPasswordErrorCode = UserErrorBase + 2 // Incorrect password
	UserNotVerifyCode = UserErrorBase + 3 // User not verified
	UserLockedCode = UserErrorBase + 4 // User is locked
	UserDisabledCode = UserErrorBase + 5 // User is disabled
	UserExpiredCode = UserErrorBase + 6 // User account expired
	UserAlreadyExistsCode = UserErrorBase + 7 // User already exists
	UserNameOrPasswordErrorCode = UserErrorBase + 8 // Incorrect username or password
	UserAuthFailedCode = UserErrorBase + 9 // Authentication failed
	UserNoPermissionCode = UserErrorBase + 10 // No permission
	UserPasswordErrCode = UserErrorBase + 11 // Password error
	UserNotExistCode = UserErrorBase + 12 // User does not exist
	UserTokenErrorCode = UserErrorBase + 13 // Invalid login credentials
	UserTokenExpiredCode = UserErrorBase + 14 // Login expired, please login again
	// 权限相关错误码
	AccessDeniedCode = AuthErrorBase + 0 // Access denied
	CasbinServiceCode = AuthErrorBase + 1 // Casbin service error
	// 缓存相关错误码
	ErrorCacheCode = CacheErrorBase + 0 // Cache operation failed
	ErrorCacheTimeoutCode = CacheErrorBase + 1 // Cache operation timeout
	ErrorCacheKeyCode = CacheErrorBase + 2 // Cache key does not exist
	ErrorCacheValueCode = CacheErrorBase + 3 // Cache value error
	// 文件相关错误码
	FileStroageCode = FileErrorBase + 0 // File storage failed
)

// 预定义错误实例
var (
	// 基础错误
	ErrSuccess = NewAppError(SuccessCode, "成功", http.StatusOK).SetI18nKey("error.success")
	// 通用错误码实例
	ErrInvalidParams = NewAppError(InvalidParamsCode, "Invalid request parameters", http.StatusBadRequest).SetI18nKey("error.common.invalid_params")
	ErrUnknown = NewAppError(ErrorUnknownCode, "Server is busy, please try again later", http.StatusInternalServerError).SetI18nKey("error.common.unknown")
	ErrNotExistCert = NewAppError(ErrorNotExistCertCode, "Authentication type does not exist", http.StatusBadRequest).SetI18nKey("error.common.not_exist_cert")
	ErrNotFound = NewAppError(ErrorNotFoundCode, "Resource not found", http.StatusNotFound).SetI18nKey("error.common.not_found")
	ErrDatabase = NewAppError(ErrorDatabaseCode, "Database operation failed", http.StatusInternalServerError).SetI18nKey("error.common.database")
	ErrInternal = NewAppError(ErrorInternalCode, "Internal server error", http.StatusInternalServerError).SetI18nKey("error.common.internal")
	ErrError = NewAppError(ErrorCode, "Error", http.StatusInternalServerError).SetI18nKey("error.common.")
	ErrParam = NewAppError(ErrorParamCode, "Parameter error", http.StatusBadRequest).SetI18nKey("error.common.param")
	// 数据库错误码实例
	ErrDatabaseInsertError = NewAppError(DatabaseInsertErrorCode, "Database insert failed", http.StatusInternalServerError).SetI18nKey("error.database.database_insert_error")
	ErrDatabaseDeleteError = NewAppError(DatabaseDeleteErrorCode, "Database delete failed", http.StatusInternalServerError).SetI18nKey("error.database.database_delete_error")
	ErrDatabaseQueryError = NewAppError(DatabaseQueryErrorCode, "Database query failed", http.StatusInternalServerError).SetI18nKey("error.database.database_query_error")
	// 用户相关错误码实例
	ErrUserNoLogin = NewAppError(UserNoLoginCode, "User not logged in", http.StatusUnauthorized).SetI18nKey("error.user.user_no_login")
	ErrUserNotFound = NewAppError(UserNotFoundCode, "User not found", http.StatusNotFound).SetI18nKey("error.user.user_not_found")
	ErrUserPasswordError = NewAppError(UserPasswordErrorCode, "Incorrect password", http.StatusUnauthorized).SetI18nKey("error.user.user_password_error")
	ErrUserNotVerify = NewAppError(UserNotVerifyCode, "User not verified", http.StatusUnauthorized).SetI18nKey("error.user.user_not_verify")
	ErrUserLocked = NewAppError(UserLockedCode, "User is locked", http.StatusUnauthorized).SetI18nKey("error.user.user_locked")
	ErrUserDisabled = NewAppError(UserDisabledCode, "User is disabled", http.StatusUnauthorized).SetI18nKey("error.user.user_disabled")
	ErrUserExpired = NewAppError(UserExpiredCode, "User account expired", http.StatusUnauthorized).SetI18nKey("error.user.user_expired")
	ErrUserAlreadyExists = NewAppError(UserAlreadyExistsCode, "User already exists", http.StatusUnauthorized).SetI18nKey("error.user.user_already_exists")
	ErrUserNameOrPasswordError = NewAppError(UserNameOrPasswordErrorCode, "Incorrect username or password", http.StatusUnauthorized).SetI18nKey("error.user.user_name_or_password_error")
	ErrUserAuthFailed = NewAppError(UserAuthFailedCode, "Authentication failed", http.StatusUnauthorized).SetI18nKey("error.user.user_auth_failed")
	ErrUserNoPermission = NewAppError(UserNoPermissionCode, "No permission", http.StatusUnauthorized).SetI18nKey("error.user.user_no_permission")
	ErrUserPasswordErr = NewAppError(UserPasswordErrCode, "Password error", http.StatusUnauthorized).SetI18nKey("error.user.user_password_err")
	ErrUserNotExist = NewAppError(UserNotExistCode, "User does not exist", http.StatusUnauthorized).SetI18nKey("error.user.user_not_exist")
	ErrUserTokenError = NewAppError(UserTokenErrorCode, "Invalid login credentials", http.StatusUnauthorized).SetI18nKey("error.user.user_token_error")
	ErrUserTokenExpired = NewAppError(UserTokenExpiredCode, "Login expired, please login again", http.StatusUnauthorized).SetI18nKey("error.user.user_token_expired")
	// 权限相关错误码实例
	ErrAccessDenied = NewAppError(AccessDeniedCode, "Access denied", http.StatusForbidden).SetI18nKey("error.auth.access_denied")
	ErrCasbinService = NewAppError(CasbinServiceCode, "Casbin service error", http.StatusInternalServerError).SetI18nKey("error.auth.casbin_service")
	// 缓存相关错误码实例
	ErrCache = NewAppError(ErrorCacheCode, "Cache operation failed", http.StatusInternalServerError).SetI18nKey("error.cache.cache")
	ErrCacheTimeout = NewAppError(ErrorCacheTimeoutCode, "Cache operation timeout", http.StatusInternalServerError).SetI18nKey("error.cache.cache_timeout")
	ErrCacheKey = NewAppError(ErrorCacheKeyCode, "Cache key does not exist", http.StatusInternalServerError).SetI18nKey("error.cache.cache_key")
	ErrCacheValue = NewAppError(ErrorCacheValueCode, "Cache value error", http.StatusInternalServerError).SetI18nKey("error.cache.cache_value")
	// 文件相关错误码实例
	ErrFileStroage = NewAppError(FileStroageCode, "File storage failed", http.StatusInternalServerError).SetI18nKey("error.file.file_stroage")
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