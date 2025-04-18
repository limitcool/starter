package errorx

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// 数据库错误类型
const (
	DBErrorTypeNotFound    = "not_found"
	DBErrorTypeDuplicate   = "duplicate"
	DBErrorTypeForeignKey  = "foreign_key"
	DBErrorTypeConstraint  = "constraint"
	DBErrorTypeConnection  = "connection"
	DBErrorTypeTransaction = "transaction"
	DBErrorTypeUnknown     = "unknown"
)

// HandleDBError 处理数据库错误
func HandleDBError(err error) error {
	if err == nil {
		return nil
	}

	// 处理常见的数据库错误
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound), errors.Is(err, sql.ErrNoRows):
		return ErrNotFound.WithError(err)
	case isDuplicateError(err):
		tmpErr := ErrDatabaseInsertError.WithMsg("记录已存在")
		return tmpErr.WithError(err)
	case isForeignKeyError(err):
		tmpErr := ErrDatabaseInsertError.WithMsg("外键约束失败")
		return tmpErr.WithError(err)
	case isConstraintError(err):
		tmpErr := ErrDatabaseInsertError.WithMsg("约束检查失败")
		return tmpErr.WithError(err)
	case isConnectionError(err):
		tmpErr := ErrDatabase.WithMsg("数据库连接失败")
		return tmpErr.WithError(err)
	case isTransactionError(err):
		tmpErr := ErrDatabase.WithMsg("事务处理失败")
		return tmpErr.WithError(err)
	default:
		return ErrDatabase.WithError(err)
	}
}

// GetDBErrorType 获取数据库错误类型
func GetDBErrorType(err error) string {
	if err == nil {
		return ""
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound), errors.Is(err, sql.ErrNoRows):
		return DBErrorTypeNotFound
	case isDuplicateError(err):
		return DBErrorTypeDuplicate
	case isForeignKeyError(err):
		return DBErrorTypeForeignKey
	case isConstraintError(err):
		return DBErrorTypeConstraint
	case isConnectionError(err):
		return DBErrorTypeConnection
	case isTransactionError(err):
		return DBErrorTypeTransaction
	default:
		return DBErrorTypeUnknown
	}
}

// isDuplicateError 判断是否为重复键错误
func isDuplicateError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "duplicate") ||
		strings.Contains(errMsg, "Duplicate") ||
		strings.Contains(errMsg, "unique constraint") ||
		strings.Contains(errMsg, "UNIQUE constraint")
}

// isForeignKeyError 判断是否为外键错误
func isForeignKeyError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "foreign key") ||
		strings.Contains(errMsg, "FOREIGN KEY")
}

// isConstraintError 判断是否为约束错误
func isConstraintError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "constraint") ||
		strings.Contains(errMsg, "Constraint") ||
		strings.Contains(errMsg, "CHECK constraint")
}

// isConnectionError 判断是否为连接错误
func isConnectionError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "connection") ||
		strings.Contains(errMsg, "Connection") ||
		strings.Contains(errMsg, "dial") ||
		strings.Contains(errMsg, "timeout")
}

// isTransactionError 判断是否为事务错误
func isTransactionError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "transaction") ||
		strings.Contains(errMsg, "Transaction") ||
		strings.Contains(errMsg, "commit") ||
		strings.Contains(errMsg, "rollback")
}

// WrapIfErr 如果err不为nil，则包装为AppError
func WrapIfErr(err error, appErr *AppError) error {
	if err == nil {
		return nil
	}
	return appErr.WithError(err)
}

// WrapWithMsg 包装错误并添加自定义消息
func WrapWithMsg(err error, appErr *AppError, msg string) error {
	if err == nil {
		return nil
	}
	return appErr.WithNewMsgAndError(msg, err)
}

// FormatError 格式化错误信息
func FormatError(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
