package logger

import (
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/configs"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Setup 初始化日志配置
func Setup(config configs.LogConfig) {
	var outputs []io.Writer

	// 配置日志级别
	level := parseLogLevel(config.Level)

	// 配置输出
	for _, output := range config.Output {
		switch output {
		case "console":
			outputs = append(outputs, os.Stdout)
		case "file":
			outputs = append(outputs, &lumberjack.Logger{
				Filename:   config.FileConfig.Path,
				MaxSize:    config.FileConfig.MaxSize, // MB
				MaxAge:     config.FileConfig.MaxAge,  // days
				MaxBackups: config.FileConfig.MaxBackups,
				Compress:   config.FileConfig.Compress,
			})
		}
	}

	// 如果没有配置输出，默认输出到控制台
	if len(outputs) == 0 {
		outputs = append(outputs, os.Stdout)
	}

	// 创建多输出writer
	multiWriter := io.MultiWriter(outputs...)

	// 确定日志格式
	formatter := parseLogFormat(config.Format)

	// 配置全局logger
	log.SetDefault(log.NewWithOptions(multiWriter, log.Options{
		Level:           level,
		Prefix:          "🌏 starter",
		TimeFormat:      time.RFC3339,
		ReportTimestamp: true,
		ReportCaller:    level == log.DebugLevel,
		Formatter:       formatter,
	}))
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
