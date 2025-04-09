package errorx

// 成功码
const (
	SuccessCode = 0
)

// 通用错误码 (10000-19999)
const (
	InvalidParamsCode = 10000 + iota
	ErrorUnknownCode
	ErrorNotExistCertCode
	ErrorNotFoundCode // 资源不存在
	ErrorDatabaseCode // 数据库操作失败
	ErrorInternalCode // 服务器内部错误
	ErrorCode         // 通用错误
	ErrorParamCode    // 参数错误
)

// 数据库错误码 (20000-29999)
const (
	DatabaseInsertErrorCode = 20000 + iota // 数据库插入错误
	DatabaseDeleteErrorCode                // 数据库删除失败
	DatabaseQueryErrorCode                 // 数据库查询失败
)

// 用户相关错误码 (30000-39999)
const (
	UserNoLoginCode             = 30000 + iota // 用户未登录
	UserNotFoundCode                           // 用户不存在
	UserPasswordErrorCode                      // 用户密码错误
	UserNotVerifyCode                          // 用户未验证
	UserLockedCode                             // 用户被锁定
	UserDisabledCode                           // 用户被关闭
	UserExpiredCode                            // 用户过期
	UserAlreadyExistsCode                      // 用户已存在
	UserNameOrPasswordErrorCode                // 用户名或密码错误
	UserAuthFailedCode                         // 用户鉴权失败
	UserNoPermissionCode                       // 用户无权访问
	UserPasswordErrCode                        // 密码错误(兼容旧版)
	UserNotExistCode                           // 用户不存在(兼容旧版)
	UserTokenErrorCode                         // 用户Token错误
	UserTokenExpiredCode                       // 用户Token过期
)

// 权限相关错误码 (40000-49999)
const (
	AccessDeniedCode  = 40000 + iota // 访问被拒绝
	CasbinServiceCode                // Casbin服务错误
)

// 缓存相关错误码 (50000-59999)
const (
	ErrorCacheCode        = 50000 + iota // 缓存操作失败
	ErrorCacheTimeoutCode                // 缓存操作超时
	ErrorCacheKeyCode                    // 缓存键不存在
	ErrorCacheValueCode                  // 缓存值错误
)

const (
	FileStroageCode = 60000 + iota // 文件存储失败
)

// 预定义错误消息，用于初始化AppError实例
var (
	// 通用错误消息
	successMsg           = "成功"
	invalidParamsMsg     = "请求参数错误"
	errorUnknownMsg      = "服务器开小差啦，稍后再来试一试"
	errorNotExistCertMsg = "不存在的认证类型"
	errorNotFoundMsg     = "资源不存在"
	errorDatabaseMsg     = "数据库操作失败"
	errorInternalMsg     = "服务器内部错误"
	errorMsg             = "错误"
	errorParamMsg        = "参数错误"

	// 数据库错误消息
	databaseInsertErrorMsg = "数据库插入失败"
	databaseDeleteErrorMsg = "数据库删除失败"
	databaseQueryErrorMsg  = "数据库查询失败"

	// 用户相关错误消息
	userNoLoginMsg             = "用户未登录"
	userNotFoundMsg            = "用户不存在"
	userPasswordErrorMsg       = "密码错误"
	userNotVerifyMsg           = "用户未验证"
	userLockedMsg              = "用户已锁定"
	userDisabledMsg            = "用户已禁用"
	userExpiredMsg             = "用户已过期"
	userAlreadyExistsMsg       = "用户已存在"
	userNameOrPasswordErrorMsg = "用户名或密码错误"
	userAuthFailedMsg          = "认证失败"
	userNoPermissionMsg        = "没有权限"
	userPasswordErrMsg         = "密码错误"
	userNotExistMsg            = "用户不存在"
	userTokenErrorMsg          = "登录凭证无效"
	userTokenExpiredMsg        = "登录已过期，请重新登录"

	// 权限相关错误消息
	accessDeniedMsg  = "访问被拒绝"
	casbinServiceMsg = "Casbin服务错误"

	// 缓存相关错误消息
	errorCacheMsg        = "缓存操作失败"
	errorCacheTimeoutMsg = "缓存操作超时"
	errorCacheKeyMsg     = "缓存键不存在"
	errorCacheValueMsg   = "缓存值错误"

	// 文件相关错误消息
	fileStroageMsg = "文件存储失败"
)

var (
	ErrSuccess = NewErrCodeMsg(SuccessCode, successMsg)
)

// 通用错误实例
var (
	ErrInvalidParams = NewErrCodeMsg(InvalidParamsCode, invalidParamsMsg)
	ErrUnknown       = NewErrCodeMsg(ErrorUnknownCode, errorUnknownMsg)
	ErrNotExistCert  = NewErrCodeMsg(ErrorNotExistCertCode, errorNotExistCertMsg)
	ErrNotFound      = NewErrCodeMsg(ErrorNotFoundCode, errorNotFoundMsg)
	ErrDatabase      = NewErrCodeMsg(ErrorDatabaseCode, errorDatabaseMsg)
	ErrInternal      = NewErrCodeMsg(ErrorInternalCode, errorInternalMsg)
	ErrGeneral       = NewErrCodeMsg(ErrorCode, errorMsg)
	ErrParam         = NewErrCodeMsg(ErrorParamCode, errorParamMsg)
)

// 数据库错误实例
var (
	ErrDatabaseInsert = NewErrCodeMsg(DatabaseInsertErrorCode, databaseInsertErrorMsg)
	ErrDatabaseDelete = NewErrCodeMsg(DatabaseDeleteErrorCode, databaseDeleteErrorMsg)
	ErrDatabaseQuery  = NewErrCodeMsg(DatabaseQueryErrorCode, databaseQueryErrorMsg)
)

// 用户相关错误实例
var (
	ErrUserNoLogin             = NewErrCodeMsg(UserNoLoginCode, userNoLoginMsg)
	ErrUserNotFound            = NewErrCodeMsg(UserNotFoundCode, userNotFoundMsg)
	ErrUserPasswordError       = NewErrCodeMsg(UserPasswordErrorCode, userPasswordErrorMsg)
	ErrUserNotVerify           = NewErrCodeMsg(UserNotVerifyCode, userNotVerifyMsg)
	ErrUserLocked              = NewErrCodeMsg(UserLockedCode, userLockedMsg)
	ErrUserDisabled            = NewErrCodeMsg(UserDisabledCode, userDisabledMsg)
	ErrUserExpired             = NewErrCodeMsg(UserExpiredCode, userExpiredMsg)
	ErrUserAlreadyExists       = NewErrCodeMsg(UserAlreadyExistsCode, userAlreadyExistsMsg)
	ErrUserNameOrPasswordError = NewErrCodeMsg(UserNameOrPasswordErrorCode, userNameOrPasswordErrorMsg)
	ErrUserAuthFailed          = NewErrCodeMsg(UserAuthFailedCode, userAuthFailedMsg)
	ErrUserNoPermission        = NewErrCodeMsg(UserNoPermissionCode, userNoPermissionMsg)
	ErrUserPasswordErr         = NewErrCodeMsg(UserPasswordErrCode, userPasswordErrMsg)
	ErrUserNotExist            = NewErrCodeMsg(UserNotExistCode, userNotExistMsg)
	ErrUserTokenError          = NewErrCodeMsg(UserTokenErrorCode, userTokenErrorMsg)
	ErrUserTokenExpired        = NewErrCodeMsg(UserTokenExpiredCode, userTokenExpiredMsg)
)

// 权限相关错误实例
var (
	ErrAccessDenied  = NewErrCodeMsg(AccessDeniedCode, accessDeniedMsg)
	ErrCasbinService = NewErrCodeMsg(CasbinServiceCode, casbinServiceMsg)
)

// 缓存相关错误实例
var (
	ErrCache        = NewErrCodeMsg(ErrorCacheCode, errorCacheMsg)
	ErrCacheTimeout = NewErrCodeMsg(ErrorCacheTimeoutCode, errorCacheTimeoutMsg)
	ErrCacheKey     = NewErrCodeMsg(ErrorCacheKeyCode, errorCacheKeyMsg)
	ErrCacheValue   = NewErrCodeMsg(ErrorCacheValueCode, errorCacheValueMsg)
)

// 文件相关错误实例
var (
	ErrFileStroage = NewErrCodeMsg(FileStroageCode, fileStroageMsg)
)
