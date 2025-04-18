package errorx

// 错误类型常量
const (
	// 错误类型
	ErrorTypeUnknown     = "unknown"     // 未知错误
	ErrorTypeValidation  = "validation"  // 验证错误
	ErrorTypeDatabase    = "database"    // 数据库错误
	ErrorTypeAuth        = "auth"        // 认证错误
	ErrorTypePermission  = "permission"  // 权限错误
	ErrorTypeNotFound    = "not_found"   // 资源未找到
	ErrorTypeConflict    = "conflict"    // 资源冲突
	ErrorTypeInternal    = "internal"    // 内部错误
	ErrorTypeExternal    = "external"    // 外部服务错误
	ErrorTypeInput       = "input"       // 输入错误
	ErrorTypeOutput      = "output"      // 输出错误
	ErrorTypeTimeout     = "timeout"     // 超时错误
	ErrorTypeConnection  = "connection"  // 连接错误
	ErrorTypeFile        = "file"        // 文件错误
	ErrorTypeCache       = "cache"       // 缓存错误
	ErrorTypeConfig      = "config"      // 配置错误
	ErrorTypeMiddleware  = "middleware"  // 中间件错误
	ErrorTypeController  = "controller"  // 控制器错误
	ErrorTypeService     = "service"     // 服务错误
	ErrorTypeRepository  = "repository"  // 仓库错误
	ErrorTypeModel       = "model"       // 模型错误
	ErrorTypeAPI         = "api"         // API错误
	ErrorTypeHTTP        = "http"        // HTTP错误
	ErrorTypeGRPC        = "grpc"        // gRPC错误
	ErrorTypeWebsocket   = "websocket"   // WebSocket错误
	ErrorTypeJSON        = "json"        // JSON错误
	ErrorTypeXML         = "xml"         // XML错误
	ErrorTypeYAML        = "yaml"        // YAML错误
	ErrorTypeEncoding    = "encoding"    // 编码错误
	ErrorTypeDecoding    = "decoding"    // 解码错误
	ErrorTypeSecurity    = "security"    // 安全错误
	ErrorTypeRate        = "rate"        // 速率错误
	ErrorTypeQuota       = "quota"       // 配额错误
	ErrorTypeLimit       = "limit"       // 限制错误
	ErrorTypeThrottle    = "throttle"    // 节流错误
	ErrorTypeCircuit     = "circuit"     // 熔断错误
	ErrorTypeFallback    = "fallback"    // 降级错误
	ErrorTypeRetry       = "retry"       // 重试错误
	ErrorTypeBackoff     = "backoff"     // 退避错误
	ErrorTypeDeadline    = "deadline"    // 截止错误
	ErrorTypeCancel      = "cancel"      // 取消错误
	ErrorTypeInterrupt   = "interrupt"   // 中断错误
	ErrorTypeAbort       = "abort"       // 中止错误
	ErrorTypePanic       = "panic"       // 恐慌错误
	ErrorTypeRecover     = "recover"     // 恢复错误
)

// GetErrorType 获取错误类型
func GetErrorType(err error) string {
	if err == nil {
		return ""
	}

	// 如果是AppError类型，根据错误码判断类型
	if appErr, ok := err.(*AppError); ok {
		code := appErr.GetErrorCode()
		switch {
		case code >= CommonErrorBase && code < DatabaseErrorBase:
			return ErrorTypeUnknown
		case code >= DatabaseErrorBase && code < UserErrorBase:
			return ErrorTypeDatabase
		case code >= UserErrorBase && code < AuthErrorBase:
			return ErrorTypeAuth
		case code >= AuthErrorBase && code < CacheErrorBase:
			return ErrorTypePermission
		case code >= CacheErrorBase && code < FileErrorBase:
			return ErrorTypeCache
		case code >= FileErrorBase && code < HttpErrorBase:
			return ErrorTypeFile
		case code >= HttpErrorBase:
			return ErrorTypeHTTP
		default:
			return ErrorTypeUnknown
		}
	}

	// 对于非AppError类型，返回未知类型
	return ErrorTypeUnknown
}

// IsErrorType 判断错误是否为指定类型
func IsErrorType(err error, errorType string) bool {
	return GetErrorType(err) == errorType
}
