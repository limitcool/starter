package logger

import (
	"sync"

	"github.com/charmbracelet/log"
)

// 日志堆栈相关配置
type StackTraceConfig struct {
	Enabled        bool   // 是否启用堆栈跟踪
	Level          string // 记录堆栈的最低日志级别
	MaxStackFrames int    // 堆栈帧最大数量
}

var (
	// 堆栈跟踪配置
	stackTraceConfig StackTraceConfig
	// 配置锁
	configLock sync.RWMutex
)

// UpdateStackTraceConfig 更新堆栈跟踪配置
func UpdateStackTraceConfig(enabled bool, level string, maxFrames int) {
	configLock.Lock()
	defer configLock.Unlock()

	// 设置默认值
	if maxFrames <= 0 {
		maxFrames = 20 // 默认最多显示20帧
	}

	// 如果未指定堆栈跟踪级别，默认为error
	if level == "" {
		level = "error"
	}

	stackTraceConfig = StackTraceConfig{
		Enabled:        enabled,
		Level:          level,
		MaxStackFrames: maxFrames,
	}
}

// GetStackTraceConfig 获取当前堆栈跟踪配置
func GetStackTraceConfig() StackTraceConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return stackTraceConfig
}

// ShouldShowStackTrace 判断是否应该显示堆栈跟踪
func ShouldShowStackTrace(level log.Level) bool {
	config := GetStackTraceConfig()
	if !config.Enabled {
		return false
	}

	// 根据配置的级别决定是否显示堆栈
	configLevel := parseLogLevel(config.Level)
	return level >= configLevel
}

// GetMaxStackFrames 获取最大堆栈帧数
func GetMaxStackFrames() int {
	config := GetStackTraceConfig()
	if config.MaxStackFrames <= 0 {
		return 20 // 默认值
	}
	return config.MaxStackFrames
}
