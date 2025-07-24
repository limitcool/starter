// Package logger 提供统一的日志接口和实现
package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/pkg/errors"
)

// Level 日志级别
type Level int

const (
	// DebugLevel 调试级别
	DebugLevel Level = iota
	// InfoLevel 信息级别
	InfoLevel
	// WarnLevel 警告级别
	WarnLevel
	// ErrorLevel 错误级别
	ErrorLevel
	// FatalLevel 致命错误级别
	FatalLevel
)

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

// Format 日志格式
type Format int

const (
	// TextFormat 文本格式
	TextFormat Format = iota
	// JSONFormat JSON格式
	JSONFormat
)

// String 返回日志格式的字符串表示
func (f Format) String() string {
	switch f {
	case TextFormat:
		return "text"
	case JSONFormat:
		return "json"
	default:
		return "unknown"
	}
}

// Logger 日志接口
type Logger interface {
	// Debug 记录调试级别日志
	Debug(msg string, keysAndValues ...any)
	// Info 记录信息级别日志
	Info(msg string, keysAndValues ...any)
	// Warn 记录警告级别日志
	Warn(msg string, keysAndValues ...any)
	// Error 记录错误级别日志
	Error(msg string, keysAndValues ...any)
	// Fatal 记录致命错误级别日志并退出程序
	Fatal(msg string, keysAndValues ...any)

	// DebugContext 使用上下文记录调试级别日志
	DebugContext(ctx context.Context, msg string, keysAndValues ...any)
	// InfoContext 使用上下文记录信息级别日志
	InfoContext(ctx context.Context, msg string, keysAndValues ...any)
	// WarnContext 使用上下文记录警告级别日志
	WarnContext(ctx context.Context, msg string, keysAndValues ...any)
	// ErrorContext 使用上下文记录错误级别日志
	ErrorContext(ctx context.Context, msg string, keysAndValues ...any)
	// FatalContext 使用上下文记录致命错误级别日志并退出程序
	FatalContext(ctx context.Context, msg string, keysAndValues ...any)

	// WithFields 创建一个带有字段的新日志记录器
	WithFields(fields map[string]any) Logger
	// WithField 创建一个带有单个字段的新日志记录器
	WithField(key string, value any) Logger
	// WithContext 创建一个带有上下文的新日志记录器
	WithContext(ctx context.Context) Logger

	// SetLevel 设置日志级别
	SetLevel(level Level)
	// GetLevel 获取日志级别
	GetLevel() Level

	// SetOutput 设置日志输出
	SetOutput(w io.Writer)
	// SetFormat 设置日志格式
	SetFormat(format Format)
}

// global 全局日志记录器
var global Logger

// Default 获取默认日志记录器
func Default() Logger {
	if global == nil {
		global = NewZapLogger(nil, InfoLevel, TextFormat)
	}
	return global
}

// SetDefault 设置默认日志记录器
func SetDefault(logger Logger) {
	global = logger
}

// Debug 使用默认日志记录器记录调试级别日志
func Debug(msg string, keysAndValues ...any) {
	Default().Debug(msg, keysAndValues...)
}

// Info 使用默认日志记录器记录信息级别日志
func Info(msg string, keysAndValues ...any) {
	Default().Info(msg, keysAndValues...)
}

// Warn 使用默认日志记录器记录警告级别日志
func Warn(msg string, keysAndValues ...any) {
	Default().Warn(msg, keysAndValues...)
}

// Error 使用默认日志记录器记录错误级别日志
func Error(msg string, keysAndValues ...any) {
	Default().Error(msg, keysAndValues...)
}

// Fatal 使用默认日志记录器记录致命错误级别日志并退出程序
func Fatal(msg string, keysAndValues ...any) {
	Default().Fatal(msg, keysAndValues...)
}

// WithFields 使用默认日志记录器创建一个带有字段的新日志记录器
func WithFields(fields map[string]any) Logger {
	return Default().WithFields(fields)
}

// WithField 使用默认日志记录器创建一个带有单个字段的新日志记录器
func WithField(key string, value any) Logger {
	return Default().WithField(key, value)
}

// LogError 记录错误并添加上下文信息
func LogError(msg string, err error, keysAndValues ...any) {
	LogErrorContext(context.Background(), msg, err, keysAndValues...)
}

// LogErrorContext 使用上下文记录错误并添加上下文信息
func LogErrorContext(ctx context.Context, msg string, err error, keysAndValues ...any) {
	fields := make([]any, 0, len(keysAndValues)+2)
	fields = append(fields, "error", err)
	fields = append(fields, keysAndValues...)
	ErrorContext(ctx, msg, fields...)
}

// LogWarn 记录警告并添加上下文信息
func LogWarn(msg string, err error, keysAndValues ...any) {
	LogWarnContext(context.Background(), msg, err, keysAndValues...)
}

// LogWarnContext 使用上下文记录警告并添加上下文信息
func LogWarnContext(ctx context.Context, msg string, err error, keysAndValues ...any) {
	fields := make([]any, 0, len(keysAndValues)+2)
	fields = append(fields, "error", err)
	fields = append(fields, keysAndValues...)
	WarnContext(ctx, msg, fields...)
}

// LogInfo 记录信息并添加上下文信息
func LogInfo(msg string, keysAndValues ...any) {
	LogInfoContext(context.Background(), msg, keysAndValues...)
}

// LogInfoContext 使用上下文记录信息并添加上下文信息
func LogInfoContext(ctx context.Context, msg string, keysAndValues ...any) {
	InfoContext(ctx, msg, keysAndValues...)
}

// LogDebug 记录调试信息并添加上下文信息
func LogDebug(msg string, keysAndValues ...any) {
	LogDebugContext(context.Background(), msg, keysAndValues...)
}

// LogDebugContext 使用上下文记录调试信息并添加上下文信息
func LogDebugContext(ctx context.Context, msg string, keysAndValues ...any) {
	DebugContext(ctx, msg, keysAndValues...)
}

// DebugContext 使用上下文记录调试级别日志
func DebugContext(ctx context.Context, msg string, keysAndValues ...any) {
	Default().DebugContext(ctx, msg, keysAndValues...)
}

// InfoContext 使用上下文记录信息级别日志
func InfoContext(ctx context.Context, msg string, keysAndValues ...any) {
	Default().InfoContext(ctx, msg, keysAndValues...)
}

// WarnContext 使用上下文记录警告级别日志
func WarnContext(ctx context.Context, msg string, keysAndValues ...any) {
	Default().WarnContext(ctx, msg, keysAndValues...)
}

// ErrorContext 使用上下文记录错误级别日志
func ErrorContext(ctx context.Context, msg string, keysAndValues ...any) {
	Default().ErrorContext(ctx, msg, keysAndValues...)
}

// FatalContext 使用上下文记录致命错误级别日志并退出程序
func FatalContext(ctx context.Context, msg string, keysAndValues ...any) {
	Default().FatalContext(ctx, msg, keysAndValues...)
}

// WithContext 使用默认日志记录器创建一个带有上下文的新日志记录器
func WithContext(ctx context.Context) Logger {
	return Default().WithContext(ctx)
}

// LogErrorWithContext 记录错误日志，包含错误详情和堆栈信息，并使用上下文
func LogErrorWithContext(ctx context.Context, msg string, err error, keyvals ...any) {
	// 构建日志字段
	fields := make([]any, 0, len(keyvals)+4) // 预分配空间

	// 检查是否需要显示堆栈
	showStackTrace := ShouldShowStackTrace(ErrorLevel)

	// 判断错误类型并处理
	if err != nil {
		// 添加基本错误信息
		fields = append(fields, "error", err.Error())

		// 处理 AppError 类型
		if appErr, ok := err.(*errorx.AppError); ok {
			// 添加错误码
			fields = append(fields, "error_code", appErr.Code())

			// 添加错误链
			errorChain := fmt.Sprintf("%+v", appErr)
			if errorChain != "" {
				fields = append(fields, "error_chain", errorChain)

				// 如果需要显示堆栈，并且错误中没有包含堆栈信息，则尝试添加
				if showStackTrace && !strings.Contains(errorChain, "[") {
					// 尝试获取和添加堆栈信息
					if formatter, ok := err.(fmt.Formatter); ok {
						var buf bytes.Buffer
						_, _ = fmt.Fprintf(&buf, "%+v", formatter)
						fields = append(fields, "stack_trace", "\n"+buf.String())
					}
				}
			}
		} else {
			// 非 AppError 类型，如果需要显示堆栈，则尝试添加
			if showStackTrace {
				// 检查是否是 pkg/errors 类型的错误
				var pkgErr interface{ StackTrace() errors.StackTrace }
				if errors.As(err, &pkgErr) {
					var buf bytes.Buffer
					_, _ = fmt.Fprintf(&buf, "%+v", pkgErr.StackTrace())
					fields = append(fields, "stack_trace", "\n"+buf.String())
				} else {
					// 尝试获取和添加堆栈信息
					if formatter, ok := err.(fmt.Formatter); ok {
						var buf bytes.Buffer
						_, _ = fmt.Fprintf(&buf, "%+v", formatter)
						fields = append(fields, "stack_trace", "\n"+buf.String())
					}
				}
			}
		}
	}

	// 添加额外的字段
	fields = append(fields, keyvals...)

	// 记录错误
	ErrorContext(ctx, msg, fields...)
}
