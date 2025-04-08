package logger

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

// GinLogWriter 实现io.Writer接口，将Gin的日志输出重定向到结构化Logger
type GinLogWriter struct {
	logger *log.Logger
	level  log.Level
}

// NewGinLogWriter 创建一个新的GinLogWriter
func NewGinLogWriter(logger *log.Logger, level log.Level) *GinLogWriter {
	return &GinLogWriter{
		logger: logger,
		level:  level,
	}
}

// Write 实现io.Writer接口
func (w *GinLogWriter) Write(p []byte) (n int, err error) {
	// 将原始字符串中的非打印字符转换为其可视化表示
	message := string(p)

	// 明确捕获并替换特殊字符
	message = strings.ReplaceAll(message, "\t", "␉") // 使用特殊符号表示制表符
	message = strings.ReplaceAll(message, "\n", "␤") // 使用特殊符号表示换行符
	message = strings.TrimSuffix(message, "␤")       // 移除末尾的换行符表示

	// 处理Gin调试输出
	if strings.Contains(message, "[GIN-debug]") {
		// 对于调试消息，我们可能想要更清晰的格式
		message = strings.ReplaceAll(message, "␉", "    ") // 替换回空格以便阅读
		message = strings.ReplaceAll(message, "␤", "\n")   // 换行符保持原样
		w.logger.Debug("Gin Debug", "raw_message", message)
	} else {
		w.logWithLevel("Gin", "message", message)
	}

	return len(p), nil
}

// logWithLevel 根据日志级别记录消息
func (w *GinLogWriter) logWithLevel(msg string, keysAndValues ...interface{}) {
	switch w.level {
	case log.DebugLevel:
		w.logger.Debug(msg, keysAndValues...)
	case log.InfoLevel:
		w.logger.Info(msg, keysAndValues...)
	case log.WarnLevel:
		w.logger.Warn(msg, keysAndValues...)
	case log.ErrorLevel:
		w.logger.Error(msg, keysAndValues...)
	default:
		w.logger.Info(msg, keysAndValues...)
	}
}

// SetupGinLogger 配置Gin使用我们的结构化Logger
func SetupGinLogger() {
	// 获取默认Logger
	logger := log.Default()

	// 将Gin的标准输出重定向到我们的Logger (Debug级别)
	gin.DefaultWriter = NewGinLogWriter(logger, log.DebugLevel)

	// 将Gin的错误输出重定向到我们的Logger (Error级别)
	gin.DefaultErrorWriter = NewGinLogWriter(logger, log.ErrorLevel)
}
