package logger

import (
	"sync"

	"github.com/limitcool/starter/pkg/logconfig"
)

// 日志堆栈相关配置
type StackTraceConfig struct {
	Enabled        bool               // 是否启用堆栈跟踪
	Level          logconfig.LogLevel // 记录堆栈的最低日志级别
	MaxStackFrames int                // 堆栈帧最大数量
}

var (
	// 堆栈跟踪配置
	stackTraceConfig StackTraceConfig
	// 配置锁
	configLock sync.RWMutex
)

// UpdateStackTraceConfig 更新堆栈跟踪配置
func UpdateStackTraceConfig(enabled bool, level logconfig.LogLevel, maxFrames int) {
	configLock.Lock()
	defer configLock.Unlock()

	// 设置默认值
	if maxFrames <= 0 {
		maxFrames = 20 // 默认最多显示20帧
	}

	// 如果未指定堆栈跟踪级别，默认为error
	if level == "" || level == logconfig.LogLevel("") {
		level = logconfig.LogLevelError
	}

	stackTraceConfig = StackTraceConfig{
		Enabled:        enabled,
		Level:          level,
		MaxStackFrames: maxFrames,
	}
}

// ShouldShowStackTrace 根据日志级别判断是否应该显示堆栈跟踪
func ShouldShowStackTrace(level Level) bool {
	configLock.RLock()
	defer configLock.RUnlock()

	if !stackTraceConfig.Enabled {
		return false
	}

	// 根据配置的级别判断
	switch stackTraceConfig.Level {
	case logconfig.LogLevelDebug:
		return level >= DebugLevel
	case logconfig.LogLevelInfo:
		return level >= InfoLevel
	case logconfig.LogLevelWarn:
		return level >= WarnLevel
	case logconfig.LogLevelError:
		return level >= ErrorLevel
	case logconfig.LogLevelFatal:
		return level >= FatalLevel
	default:
		return level >= ErrorLevel
	}
}

// GetStackTraceConfig 获取堆栈跟踪配置
func GetStackTraceConfig() StackTraceConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return stackTraceConfig
}
