package logger

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/pkg/logconfig"
	"github.com/pkg/errors"
)

// Setup 初始化日志配置
func Setup(config logconfig.LogConfig) {
	// 更新堆栈跟踪配置
	UpdateStackTraceConfig(
		config.StackTraceEnabled,
		config.StackTraceLevel,
		config.MaxStackFrames,
	)

	// 创建并设置logger
	// 使用ZapLogger代替CharmLogger以提高性能
	logger := NewZapLoggerWithConfig(config)
	SetDefault(logger)
}

// parseLogLevel 解析日志级别
func parseLogLevel(level logconfig.LogLevel) Level {
	switch level {
	case logconfig.LogLevelDebug:
		return DebugLevel
	case logconfig.LogLevelInfo:
		return InfoLevel
	case logconfig.LogLevelWarn:
		return WarnLevel
	case logconfig.LogLevelError:
		return ErrorLevel
	case logconfig.LogLevelFatal:
		return FatalLevel
	default:
		return InfoLevel
	}
}

// LogErrorWithStack 记录错误日志，包含错误详情和堆栈信息
// 参数:
//   - msg: 错误消息
//   - err: 当前错误
//   - keyvals: 额外的键值对信息，按照 key1, value1, key2, value2... 格式提供
func LogErrorWithStack(msg string, err error, keyvals ...any) {
	LogErrorWithStackContext(context.Background(), msg, err, keyvals...)
}

// LogErrorWithStackContext 使用上下文记录错误日志，包含错误详情和堆栈信息
// 参数:
//   - ctx: 上下文
//   - msg: 错误消息
//   - err: 当前错误
//   - keyvals: 额外的键值对信息，按照 key1, value1, key2, value2... 格式提供
func LogErrorWithStackContext(ctx context.Context, msg string, err error, keyvals ...any) {
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
			errorChain := errorx.FormatErrorChain(appErr)
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
