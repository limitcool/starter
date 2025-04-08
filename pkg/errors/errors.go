package errors

import (
	"errors"

	"github.com/limitcool/starter/pkg/code"
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
		errors.Is(err, gorm.ErrForeignKeyViolated) {
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
		return code.ErrorNotFound, code.GetMsg(code.ErrorNotFound)
	}

	// 判断重复错误
	if IsDuplicate(err) {
		var dupErr *DuplicateError
		if errors.As(err, &dupErr) && dupErr.Resource != "" {
			return code.UserAlreadyExists, dupErr.Resource + "已存在"
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

	// 其他错误
	return code.ErrorUnknown, code.GetMsg(code.ErrorUnknown)
}
