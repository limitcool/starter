// Package logconfig 提供日志配置结构体和常量
package logconfig

import "time"

// LogFormat 日志格式类型
type LogFormat string

// 日志格式常量
const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

// LogStyle 日志风格类型
type LogStyle string

// 日志风格常量
const (
	LogStylePlain      LogStyle = "plain"      // 非结构化日志，使用键值对
	LogStyleStructured LogStyle = "structured" // 结构化日志，使用对象
)

// LogLevel 日志级别类型
type LogLevel string

// 日志级别常量
const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// LogConfig 日志配置
type LogConfig struct {
	Level             LogLevel      `yaml:"level" json:"level"`                             // 日志级别
	Format            LogFormat     `yaml:"format" json:"format"`                           // 日志格式
	Style             LogStyle      `yaml:"style" json:"style"`                             // 日志风格（结构化或非结构化）
	Output            []string      `yaml:"output" json:"output"`                           // 日志输出位置
	FileConfig        FileLogConfig `yaml:"file_config" json:"file_config"`                 // 文件日志配置
	StackTraceLevel   LogLevel      `yaml:"stack_trace_level" json:"stack_trace_level"`     // 堆栈跟踪级别
	StackTraceEnabled bool          `yaml:"stack_trace_enabled" json:"stack_trace_enabled"` // 是否启用堆栈跟踪
	MaxStackFrames    int           `yaml:"max_stack_frames" json:"max_stack_frames"`       // 最大堆栈帧数
	Sampling          bool          `yaml:"sampling" json:"sampling"`                       // 是否启用采样（高频日志降频）
	Development       bool          `yaml:"development" json:"development"`                 // 是否为开发模式（更详细的日志）
	EncoderConfig     EncoderConfig `yaml:"encoder_config" json:"encoder_config"`           // 编码器配置
}

// FileLogConfig 文件日志配置
type FileLogConfig struct {
	Path       string        `yaml:"path" json:"path"`               // 日志文件路径
	MaxSize    int           `yaml:"max_size" json:"max_size"`       // 单个日志文件最大大小(MB)
	MaxAge     int           `yaml:"max_age" json:"max_age"`         // 日志文件保留天数
	MaxBackups int           `yaml:"max_backups" json:"max_backups"` // 最大备份数
	Compress   bool          `yaml:"compress" json:"compress"`       // 是否压缩
	Rotation   time.Duration `yaml:"rotation" json:"rotation"`       // 日志轮转时间间隔
}

// EncoderConfig 编码器配置
type EncoderConfig struct {
	MessageKey     string `yaml:"message_key" json:"message_key"`         // 消息字段名
	LevelKey       string `yaml:"level_key" json:"level_key"`             // 级别字段名
	TimeKey        string `yaml:"time_key" json:"time_key"`               // 时间字段名
	NameKey        string `yaml:"name_key" json:"name_key"`               // 名称字段名
	CallerKey      string `yaml:"caller_key" json:"caller_key"`           // 调用者字段名
	StacktraceKey  string `yaml:"stacktrace_key" json:"stacktrace_key"`   // 堆栈字段名
	EncodeTime     string `yaml:"encode_time" json:"encode_time"`         // 时间编码方式
	EncodeLevel    string `yaml:"encode_level" json:"encode_level"`       // 级别编码方式
	EncodeCaller   string `yaml:"encode_caller" json:"encode_caller"`     // 调用者编码方式
	EncodeDuration string `yaml:"encode_duration" json:"encode_duration"` // 时间间隔编码方式
}

// DefaultEncoderConfig 返回默认编码器配置
func DefaultEncoderConfig() EncoderConfig {
	return EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		EncodeTime:     "iso8601",
		EncodeLevel:    "capital",
		EncodeCaller:   "short",
		EncodeDuration: "string",
	}
}

// DefaultLogConfig 返回默认日志配置
func DefaultLogConfig() LogConfig {
	return LogConfig{
		Level:  LogLevelInfo,
		Format: LogFormatText,
		Style:  LogStylePlain,
		Output: []string{"console"},
		FileConfig: FileLogConfig{
			Path:       "logs",
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 10,
			Compress:   false,
		},
		StackTraceEnabled: true,
		StackTraceLevel:   LogLevelError,
		MaxStackFrames:    64,
		Sampling:          false,
		Development:       false,
		EncoderConfig:     DefaultEncoderConfig(),
	}
}
