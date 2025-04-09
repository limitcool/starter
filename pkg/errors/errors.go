package errors

import (
	"errors"
	"strings"

	"github.com/limitcool/starter/pkg/code"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// 标准错误类型
var (
	// 包装golang标准errors包
	New    = errors.New
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)

// 定义自定义错误类型

// NotFoundError 表示资源不存在的错误
type NotFoundError struct {
	Resource string
	ID       string
	Err      error
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return "resource " + e.Resource + " with id " + e.ID + " not found"
	}
	return "resource " + e.Resource + " not found"
}

func (e *NotFoundError) Unwrap() error { return e.Err }

// NewNotFoundError 创建一个NotFoundError
func NewNotFoundError(resource, id string) error {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// DuplicateError 表示资源重复的错误
type DuplicateError struct {
	Resource string
	Field    string
	Value    string
	Err      error
}

func (e *DuplicateError) Error() string {
	return "resource " + e.Resource + " with " + e.Field + "=" + e.Value + " already exists"
}

func (e *DuplicateError) Unwrap() error { return e.Err }

// NewDuplicateError 创建一个DuplicateError
func NewDuplicateError(resource, field, value string) error {
	return &DuplicateError{
		Resource: resource,
		Field:    field,
		Value:    value,
	}
}

// PermissionDeniedError 表示权限不足的错误
type PermissionDeniedError struct {
	Resource string
	Action   string
	Err      error
}

func (e *PermissionDeniedError) Error() string {
	return "permission denied to " + e.Action + " " + e.Resource
}

func (e *PermissionDeniedError) Unwrap() error { return e.Err }

// NewPermissionDeniedError 创建一个PermissionDeniedError
func NewPermissionDeniedError(resource, action string) error {
	return &PermissionDeniedError{
		Resource: resource,
		Action:   action,
	}
}

// AuthenticationError 表示认证失败的错误
type AuthenticationError struct {
	Reason string
	Err    error
}

func (e *AuthenticationError) Error() string {
	if e.Reason != "" {
		return "authentication failed: " + e.Reason
	}
	return "authentication failed"
}

func (e *AuthenticationError) Unwrap() error { return e.Err }

// NewAuthenticationError 创建一个AuthenticationError
func NewAuthenticationError(reason string) error {
	return &AuthenticationError{
		Reason: reason,
	}
}

// DatabaseError 表示数据库操作错误
type DatabaseError struct {
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	if e.Operation != "" {
		return "database error during " + e.Operation + ": " + e.Err.Error()
	}
	return "database error: " + e.Err.Error()
}

func (e *DatabaseError) Unwrap() error { return e.Err }

// NewDatabaseError 创建一个DatabaseError
func NewDatabaseError(operation string, err error) error {
	return &DatabaseError{
		Operation: operation,
		Err:       err,
	}
}

// CacheError 表示缓存操作错误
type CacheError struct {
	Operation string
	Key       string
	Err       error
}

func (e *CacheError) Error() string {
	if e.Key != "" {
		return "cache error during " + e.Operation + " for key " + e.Key + ": " + e.Err.Error()
	}
	return "cache error during " + e.Operation + ": " + e.Err.Error()
}

func (e *CacheError) Unwrap() error { return e.Err }

// NewCacheError 创建一个CacheError
func NewCacheError(operation, key string, err error) error {
	return &CacheError{
		Operation: operation,
		Key:       key,
		Err:       err,
	}
}

// ValidationError 表示输入验证错误
type ValidationError struct {
	Field   string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return "validation error: field " + e.Field + " " + e.Message
	}
	return "validation error: " + e.Message
}

func (e *ValidationError) Unwrap() error { return e.Err }

// NewValidationError 创建一个ValidationError
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// StorageError 表示文件存储操作错误
type StorageError struct {
	Operation string
	Path      string
	Err       error
}

func (e *StorageError) Error() string {
	if e.Path != "" {
		return "storage error during " + e.Operation + " for path " + e.Path + ": " + e.Err.Error()
	}
	return "storage error during " + e.Operation + ": " + e.Err.Error()
}

func (e *StorageError) Unwrap() error { return e.Err }

// NewStorageError 创建一个StorageError
func NewStorageError(operation, path string, err error) error {
	return &StorageError{
		Operation: operation,
		Path:      path,
		Err:       err,
	}
}

// 定义OSS存储错误码常量
const (
	ErrStorageNotFound  = "storage: file not found"
	ErrStorageExists    = "storage: file already exists"
	ErrStorageForbidden = "storage: operation forbidden"
	ErrStorageUnknown   = "storage: unknown error"
)

// WithCode 创建带错误码的错误
func WithCode(errCode int, message string) error {
	return code.NewErrCodeMsg(errCode, message)
}

// 错误类型判断函数

// IsNotFound 判断是否为"不存在"类型的错误
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	// 检查自定义NotFoundError
	var notFoundErr *NotFoundError
	if errors.As(err, &notFoundErr) {
		return true
	}

	// 检查是否为GORM的记录不存在错误
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}

	// 检查是否为自定义的NotFound错误码
	var codeErr *code.CodeError
	if errors.As(err, &codeErr) && codeErr.GetErrCode() == code.ErrorNotFound {
		return true
	}

	return false
}

// IsDuplicate 判断是否为"重复"类型的错误
func IsDuplicate(err error) bool {
	if err == nil {
		return false
	}

	// 检查自定义DuplicateError
	var dupErr *DuplicateError
	if errors.As(err, &dupErr) {
		return true
	}

	// 检查是否为GORM的重复键错误
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// 检查是否为自定义的重复错误码
	var codeErr *code.CodeError
	if errors.As(err, &codeErr) && codeErr.GetErrCode() == code.UserAlreadyExists {
		return true
	}

	return false
}

// IsPermissionDenied 判断是否为"权限拒绝"类型的错误
func IsPermissionDenied(err error) bool {
	if err == nil {
		return false
	}

	// 检查自定义PermissionDeniedError
	var permErr *PermissionDeniedError
	if errors.As(err, &permErr) {
		return true
	}

	// 检查是否为自定义的权限拒绝错误码
	var codeErr *code.CodeError
	if errors.As(err, &codeErr) && (codeErr.GetErrCode() == code.AccessDenied || codeErr.GetErrCode() == code.UserNoPermission) {
		return true
	}

	return false
}

// IsAuthenticationFailed 判断是否为"认证失败"类型的错误
func IsAuthenticationFailed(err error) bool {
	if err == nil {
		return false
	}

	// 检查自定义AuthenticationError
	var authErr *AuthenticationError
	if errors.As(err, &authErr) {
		return true
	}

	// 检查是否为自定义的认证错误码
	var codeErr *code.CodeError
	if errors.As(err, &codeErr) &&
		(codeErr.GetErrCode() == code.UserNoLogin ||
			codeErr.GetErrCode() == code.UserTokenError ||
			codeErr.GetErrCode() == code.UserTokenExpired ||
			codeErr.GetErrCode() == code.UserPasswordError ||
			codeErr.GetErrCode() == code.UserNameOrPasswordError) {
		return true
	}

	return false
}

// IsDBError 判断是否为数据库操作错误
func IsDBError(err error) bool {
	if err == nil {
		return false
	}

	// 检查自定义DatabaseError
	var dbErr *DatabaseError
	if errors.As(err, &dbErr) {
		return true
	}

	// 检查是否为GORM的数据库错误
	if errors.Is(err, gorm.ErrInvalidDB) ||
		errors.Is(err, gorm.ErrInvalidTransaction) ||
		errors.Is(err, gorm.ErrForeignKeyViolated) ||
		errors.Is(err, gorm.ErrNotImplemented) ||
		errors.Is(err, gorm.ErrMissingWhereClause) ||
		errors.Is(err, gorm.ErrUnsupportedRelation) ||
		errors.Is(err, gorm.ErrPrimaryKeyRequired) ||
		errors.Is(err, gorm.ErrModelValueRequired) ||
		errors.Is(err, gorm.ErrInvalidData) ||
		errors.Is(err, gorm.ErrUnsupportedDriver) ||
		errors.Is(err, gorm.ErrRegistered) ||
		errors.Is(err, gorm.ErrInvalidField) ||
		errors.Is(err, gorm.ErrEmptySlice) ||
		errors.Is(err, gorm.ErrDryRunModeUnsupported) {
		return true
	}

	// 检查MongoDB错误
	if isMongoDBError(err) {
		return true
	}

	// 检查是否为自定义的数据库错误码
	var codeErr *code.CodeError
	if errors.As(err, &codeErr) &&
		(codeErr.GetErrCode() == code.ErrorDatabase ||
			codeErr.GetErrCode() == code.DatabaseInsertError ||
			codeErr.GetErrCode() == code.DatabaseDeleteError ||
			codeErr.GetErrCode() == code.DatabaseQueryError) {
		return true
	}

	return false
}

// 检查是否为MongoDB错误
func isMongoDBError(err error) bool {
	if err == nil {
		return false
	}

	var mongoErr mongo.CommandError
	if errors.As(err, &mongoErr) {
		return true
	}

	var writeErr mongo.WriteException
	if errors.As(err, &writeErr) {
		return true
	}

	var bulkWriteErr mongo.BulkWriteException
	if errors.As(err, &bulkWriteErr) {
		return true
	}

	return mongo.IsTimeout(err) || mongo.IsNetworkError(err)
}

// IsCacheError 判断是否为缓存操作错误
func IsCacheError(err error) bool {
	if err == nil {
		return false
	}

	// 检查自定义CacheError
	var cacheErr *CacheError
	if errors.As(err, &cacheErr) {
		return true
	}

	// 检查是否为Redis错误
	if isRedisError(err) {
		return true
	}

	// 检查是否为自定义的缓存错误码
	var codeErr *code.CodeError
	if errors.As(err, &codeErr) &&
		(codeErr.GetErrCode() == code.ErrorCache ||
			codeErr.GetErrCode() == code.ErrorCacheTimeout ||
			codeErr.GetErrCode() == code.ErrorCacheKey ||
			codeErr.GetErrCode() == code.ErrorCacheValue) {
		return true
	}

	return false
}

// 检查是否为Redis错误
func isRedisError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	return strings.Contains(errMsg, "redis:") ||
		strings.Contains(errMsg, "connection refused") ||
		strings.Contains(errMsg, "connection timeout") ||
		strings.Contains(errMsg, "connection reset") ||
		strings.Contains(errMsg, "nil") && strings.Contains(errMsg, "redis")
}

// IsValidationError 判断是否为验证错误
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}

	// 检查自定义ValidationError
	var validErr *ValidationError
	if errors.As(err, &validErr) {
		return true
	}

	// 检查是否为自定义的参数错误码
	var codeErr *code.CodeError
	if errors.As(err, &codeErr) && codeErr.GetErrCode() == code.InvalidParams {
		return true
	}

	return false
}

// IsStorageError 判断是否为文件存储错误
func IsStorageError(err error) bool {
	if err == nil {
		return false
	}

	// 检查自定义StorageError
	var storageErr *StorageError
	if errors.As(err, &storageErr) {
		return true
	}

	// 检查是否为OSS错误
	if isOSSError(err) {
		return true
	}

	return false
}

// isOSSError 检查是否为OSS错误
func isOSSError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	return strings.Contains(errMsg, "storage:") ||
		strings.Contains(errMsg, "oss:") ||
		strings.Contains(errMsg, "s3:") ||
		strings.Contains(errMsg, "bucket") ||
		strings.Contains(errMsg, "object") ||
		strings.Contains(errMsg, "file")
}

// ParseError 解析错误并返回对应的错误码和消息
func ParseError(err error) (int, string) {
	if err == nil {
		return code.Success, code.GetMsg(code.Success)
	}

	// 如果是自定义错误码，直接使用错误码和消息
	var codeErr *code.CodeError
	if errors.As(err, &codeErr) {
		return codeErr.GetErrCode(), codeErr.GetErrMsg()
	}

	// 判断资源不存在错误
	if IsNotFound(err) {
		var notFoundErr *NotFoundError
		if errors.As(err, &notFoundErr) && notFoundErr.Resource != "" {
			return code.ErrorNotFound, notFoundErr.Resource + "不存在"
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return code.ErrorNotFound, "记录不存在"
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			return code.ErrorNotFound, "文档不存在"
		}
		return code.ErrorNotFound, code.GetMsg(code.ErrorNotFound)
	}

	// 判断重复错误
	if IsDuplicate(err) {
		var dupErr *DuplicateError
		if errors.As(err, &dupErr) && dupErr.Resource != "" {
			return code.UserAlreadyExists, dupErr.Resource + "已存在"
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return code.UserAlreadyExists, "记录已存在"
		}

		// 检查MongoDB的重复键错误
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, we := range writeErr.WriteErrors {
				if we.Code == 11000 { // MongoDB重复键错误码
					return code.UserAlreadyExists, "记录已存在"
				}
			}
		}

		return code.UserAlreadyExists, code.GetMsg(code.UserAlreadyExists)
	}

	// 判断权限错误
	if IsPermissionDenied(err) {
		var permErr *PermissionDeniedError
		if errors.As(err, &permErr) && permErr.Resource != "" && permErr.Action != "" {
			return code.AccessDenied, "您没有权限" + permErr.Action + permErr.Resource
		}
		return code.AccessDenied, code.GetMsg(code.AccessDenied)
	}

	// 判断认证错误
	if IsAuthenticationFailed(err) {
		var authErr *AuthenticationError
		if errors.As(err, &authErr) && authErr.Reason != "" {
			return code.UserAuthFailed, "认证失败：" + authErr.Reason
		}
		return code.UserAuthFailed, code.GetMsg(code.UserAuthFailed)
	}

	// 判断数据库错误
	if IsDBError(err) {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return code.ErrorDatabase, "关联数据约束错误，请先删除相关数据"
		}

		// 详细处理GORM常见错误
		if errors.Is(err, gorm.ErrInvalidTransaction) {
			return code.ErrorDatabase, "数据库事务无效"
		}
		if errors.Is(err, gorm.ErrNotImplemented) {
			return code.ErrorDatabase, "数据库操作未实现"
		}
		if errors.Is(err, gorm.ErrMissingWhereClause) {
			return code.ErrorDatabase, "数据库操作缺少WHERE条件"
		}
		if errors.Is(err, gorm.ErrUnsupportedRelation) {
			return code.ErrorDatabase, "数据库不支持的关联关系"
		}
		if errors.Is(err, gorm.ErrPrimaryKeyRequired) {
			return code.ErrorDatabase, "数据库操作需要主键"
		}
		if errors.Is(err, gorm.ErrModelValueRequired) {
			return code.ErrorDatabase, "数据库模型值为空"
		}

		// 检查MongoDB错误
		if isMongoDBError(err) {
			var cmdErr mongo.CommandError
			if errors.As(err, &cmdErr) {
				return code.ErrorDatabase, "MongoDB命令错误: " + cmdErr.Message
			}

			if mongo.IsTimeout(err) {
				return code.ErrorDatabase, "MongoDB操作超时"
			}

			if mongo.IsNetworkError(err) {
				return code.ErrorDatabase, "MongoDB网络错误"
			}

			return code.ErrorDatabase, "MongoDB操作失败"
		}

		var dbErr *DatabaseError
		if errors.As(err, &dbErr) && dbErr.Operation != "" {
			return code.ErrorDatabase, "数据库" + dbErr.Operation + "操作失败"
		}
		return code.ErrorDatabase, code.GetMsg(code.ErrorDatabase)
	}

	// 判断缓存错误
	if IsCacheError(err) {
		var cacheErr *CacheError
		if errors.As(err, &cacheErr) && cacheErr.Operation != "" {
			if cacheErr.Key != "" {
				return code.ErrorCache, "缓存" + cacheErr.Operation + "操作失败，键：" + cacheErr.Key
			}
			return code.ErrorCache, "缓存" + cacheErr.Operation + "操作失败"
		}

		// Redis错误处理
		if isRedisError(err) {
			errMsg := err.Error()
			if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "connection reset") {
				return code.ErrorCache, "Redis连接失败"
			}
			if strings.Contains(errMsg, "connection timeout") {
				return code.ErrorCacheTimeout, "Redis连接超时"
			}
			if strings.Contains(errMsg, "nil") {
				return code.ErrorCacheKey, "Redis键不存在"
			}
			return code.ErrorCache, "Redis操作失败"
		}

		return code.ErrorCache, code.GetMsg(code.ErrorCache)
	}

	// 判断验证错误
	if IsValidationError(err) {
		var validErr *ValidationError
		if errors.As(err, &validErr) {
			if validErr.Field != "" {
				return code.InvalidParams, "参数" + validErr.Field + "错误：" + validErr.Message
			}
			return code.InvalidParams, validErr.Message
		}
		return code.InvalidParams, code.GetMsg(code.InvalidParams)
	}

	// 判断文件存储错误
	if IsStorageError(err) {
		var storageErr *StorageError
		if errors.As(err, &storageErr) {
			if storageErr.Operation != "" {
				if storageErr.Path != "" {
					return code.ErrorUnknown, "文件" + storageErr.Operation + "操作失败，路径：" + storageErr.Path
				}
				return code.ErrorUnknown, "文件" + storageErr.Operation + "操作失败"
			}
		}

		// 特定OSS错误处理
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") || strings.Contains(errMsg, ErrStorageNotFound) {
			return code.ErrorNotFound, "文件不存在"
		}
		if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, ErrStorageExists) {
			return code.Error, "文件已存在"
		}
		if strings.Contains(errMsg, "forbidden") || strings.Contains(errMsg, "permission denied") || strings.Contains(errMsg, ErrStorageForbidden) {
			return code.AccessDenied, "没有文件操作权限"
		}

		return code.Error, "文件操作失败"
	}

	// 其他错误
	return code.ErrorUnknown, code.GetMsg(code.ErrorUnknown)
}
