package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

// CharmLogger 基于 charmbracelet/log 的日志实现
type CharmLogger struct {
	logger *log.Logger
	level  Level
	format Format
}

// NewCharmLogger 创建一个新的 CharmLogger
func NewCharmLogger(w io.Writer, level Level, format Format) *CharmLogger {
	if w == nil {
		w = os.Stdout
	}

	// 创建 charmbracelet/log 的 Logger
	options := log.Options{
		Level:           convertToCharmLevel(level),
		Prefix:          "🌏 starter",
		TimeFormat:      time.RFC3339,
		ReportTimestamp: true,
		ReportCaller:    level == DebugLevel,
	}

	// 设置日志格式
	if format == JSONFormat {
		options.Formatter = log.JSONFormatter
	} else {
		options.Formatter = log.TextFormatter
	}

	logger := log.NewWithOptions(w, options)

	return &CharmLogger{
		logger: logger,
		level:  level,
		format: format,
	}
}

// Debug 实现 Logger 接口
func (l *CharmLogger) Debug(msg string, keysAndValues ...any) {
	l.logger.Debug(msg, keysAndValues...)
}

// Info 实现 Logger 接口
func (l *CharmLogger) Info(msg string, keysAndValues ...any) {
	l.logger.Info(msg, keysAndValues...)
}

// Warn 实现 Logger 接口
func (l *CharmLogger) Warn(msg string, keysAndValues ...any) {
	l.logger.Warn(msg, keysAndValues...)
}

// Error 实现 Logger 接口
func (l *CharmLogger) Error(msg string, keysAndValues ...any) {
	l.logger.Error(msg, keysAndValues...)
}

// Fatal 实现 Logger 接口
func (l *CharmLogger) Fatal(msg string, keysAndValues ...any) {
	l.logger.Fatal(msg, keysAndValues...)
}

// WithFields 实现 Logger 接口
func (l *CharmLogger) WithFields(fields map[string]any) Logger {
	// 将 map 转换为 key-value 对
	keyValues := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		keyValues = append(keyValues, k, v)
	}

	// 创建新的 Logger
	newLogger := l.logger.With(keyValues...)

	return &CharmLogger{
		logger: newLogger,
		level:  l.level,
		format: l.format,
	}
}

// WithField 实现 Logger 接口
func (l *CharmLogger) WithField(key string, value any) Logger {
	newLogger := l.logger.With(key, value)

	return &CharmLogger{
		logger: newLogger,
		level:  l.level,
		format: l.format,
	}
}

// SetLevel 实现 Logger 接口
func (l *CharmLogger) SetLevel(level Level) {
	l.level = level
	l.logger.SetLevel(convertToCharmLevel(level))
}

// GetLevel 实现 Logger 接口
func (l *CharmLogger) GetLevel() Level {
	return l.level
}

// SetOutput 实现 Logger 接口
func (l *CharmLogger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

// SetFormat 实现 Logger 接口
func (l *CharmLogger) SetFormat(format Format) {
	l.format = format
	if format == JSONFormat {
		l.logger.SetFormatter(log.JSONFormatter)
	} else {
		l.logger.SetFormatter(log.TextFormatter)
	}
}

// convertToCharmLevel 将我们的日志级别转换为 charmbracelet/log 的日志级别
func convertToCharmLevel(level Level) log.Level {
	switch level {
	case DebugLevel:
		return log.DebugLevel
	case InfoLevel:
		return log.InfoLevel
	case WarnLevel:
		return log.WarnLevel
	case ErrorLevel:
		return log.ErrorLevel
	case FatalLevel:
		return log.FatalLevel
	default:
		return log.InfoLevel
	}
}

// WithContext 实现 Logger 接口
func (l *CharmLogger) WithContext(ctx context.Context) Logger {
	// 从上下文中提取关键信息
	fields := extractContextFields(ctx)

	// 将字段添加到日志中
	keyValues := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		keyValues = append(keyValues, k, v)
	}

	// 创建新的 Logger
	newLogger := l.logger.With(keyValues...)

	return &CharmLogger{
		logger: newLogger,
		level:  l.level,
		format: l.format,
	}
}

// DebugContext 实现 Logger 接口
func (l *CharmLogger) DebugContext(ctx context.Context, msg string, keysAndValues ...any) {
	l.WithContext(ctx).Debug(msg, keysAndValues...)
}

// InfoContext 实现 Logger 接口
func (l *CharmLogger) InfoContext(ctx context.Context, msg string, keysAndValues ...any) {
	l.WithContext(ctx).Info(msg, keysAndValues...)
}

// WarnContext 实现 Logger 接口
func (l *CharmLogger) WarnContext(ctx context.Context, msg string, keysAndValues ...any) {
	l.WithContext(ctx).Warn(msg, keysAndValues...)
}

// ErrorContext 实现 Logger 接口
func (l *CharmLogger) ErrorContext(ctx context.Context, msg string, keysAndValues ...any) {
	l.WithContext(ctx).Error(msg, keysAndValues...)
}

// FatalContext 实现 Logger 接口
func (l *CharmLogger) FatalContext(ctx context.Context, msg string, keysAndValues ...any) {
	l.WithContext(ctx).Fatal(msg, keysAndValues...)
}

// convertFromCharmLevel 将 charmbracelet/log 的日志级别转换为我们的日志级别
func convertFromCharmLevel(level log.Level) Level {
	switch level {
	case log.DebugLevel:
		return DebugLevel
	case log.InfoLevel:
		return InfoLevel
	case log.WarnLevel:
		return WarnLevel
	case log.ErrorLevel:
		return ErrorLevel
	case log.FatalLevel:
		return FatalLevel
	default:
		return InfoLevel
	}
}
