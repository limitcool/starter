package logger

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/configs"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Setup 初始化日志配置
func Setup(config configs.LogConfig) {
	// 配置日志级别
	level := parseLogLevel(config.Level)

	// 检查是否需要输出到控制台
	hasConsole := false
	for _, output := range config.Output {
		if output == "console" {
			hasConsole = true
			break
		}
	}

	// 如果配置为空，默认输出到控制台
	if len(config.Output) == 0 {
		hasConsole = true
	}

	// 检查是否需要输出到文件
	var fileOutput io.Writer
	for _, output := range config.Output {
		if output == "file" {
			fileOutput = &lumberjack.Logger{
				Filename:   config.FileConfig.Path,
				MaxSize:    config.FileConfig.MaxSize,
				MaxAge:     config.FileConfig.MaxAge,
				MaxBackups: config.FileConfig.MaxBackups,
				Compress:   config.FileConfig.Compress,
			}
			break
		}
	}

	// 创建基本设置
	options := log.Options{
		Level:           level,
		Prefix:          "🌏 starter",
		TimeFormat:      time.RFC3339,
		ReportTimestamp: true,
		ReportCaller:    level == log.DebugLevel,
	}

	// 根据不同情况创建logger
	var writer io.Writer

	if hasConsole && fileOutput != nil {
		// 同时输出到控制台和文件
		writer = io.MultiWriter(os.Stdout, fileOutput)
	} else if hasConsole {
		// 只输出到控制台
		writer = os.Stdout
	} else if fileOutput != nil {
		// 只输出到文件
		writer = fileOutput
	} else {
		// 默认输出到控制台
		writer = os.Stdout
	}

	// 设置日志格式
	if config.Format == configs.LogFormatJSON {
		// JSON格式
		options.Formatter = log.JSONFormatter
	} else {
		// 文本格式，支持彩色
		options.Formatter = log.TextFormatter
	}

	// 创建并设置logger
	logger := log.NewWithOptions(writer, options)
	log.SetDefault(logger)
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) log.Level {
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}

// parseLogFormat 根据配置解析日志格式
func parseLogFormat(format configs.LogFormat) log.Formatter {
	switch format {
	case configs.LogFormatJSON:
		return log.JSONFormatter
	case configs.LogFormatText:
		return log.TextFormatter
	default:
		// 默认使用文本格式
		return log.TextFormatter
	}
}

// LogError 记录错误信息，同时记录原始错误
// 参数:
//   - msg: 错误消息
//   - err: 当前错误
//   - keyvals: 额外的键值对信息，按照 key1, value1, key2, value2... 格式提供
func LogError(msg string, err error, keyvals ...interface{}) {
	// 构建日志字段
	fields := []interface{}{"err", err}

	// 添加原始错误（如果存在）
	if originalErr := errors.Unwrap(err); originalErr != nil {
		fields = append(fields, "original_err", originalErr)
	}

	// 添加额外的字段
	fields = append(fields, keyvals...)

	// 记录错误
	log.Error(msg, fields...)
}
