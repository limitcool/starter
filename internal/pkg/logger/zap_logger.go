package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/limitcool/starter/pkg/logconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger 基于 zap 的日志实现
type ZapLogger struct {
	logger        *zap.SugaredLogger
	structLogger  *zap.Logger // 结构化日志器
	level         Level
	format        Format
	style         logconfig.LogStyle
	development   bool
	sampling      bool
	encoderConfig logconfig.EncoderConfig
}

// NewZapLogger 创建一个新的 ZapLogger
func NewZapLogger(w io.Writer, level Level, format Format) *ZapLogger {
	return newZapLogger(w, level, format, nil)
}

// NewZapLoggerWithConfig 使用配置创建一个新的 ZapLogger
func NewZapLoggerWithConfig(config logconfig.LogConfig) *ZapLogger {
	// 解析日志级别
	level := parseLogLevel(config.Level)

	// 设置堆栈跟踪级别
	stackLevel := zapcore.ErrorLevel
	if config.StackTraceEnabled {
		stackLevel = convertToZapLevel(parseLogLevel(config.StackTraceLevel))
	} else {
		// 如果禁用堆栈跟踪，设置为一个非常高的级别
		stackLevel = zapcore.FatalLevel + 1
	}

	// 创建多个core，支持不同格式的输出
	var cores []zapcore.Core

	// 检查是否需要输出到控制台
	hasConsole := false
	hasFile := false
	for _, output := range config.Output {
		if output == "console" {
			hasConsole = true
		} else if output == "file" {
			hasFile = true
		}
	}

	// 如果配置为空，默认输出到控制台
	if len(config.Output) == 0 {
		hasConsole = true
	}

	// 添加控制台输出 - 使用text格式
	if hasConsole {
		consoleCore := createCore(os.Stdout, level, TextFormat, config)
		cores = append(cores, consoleCore)
	}

	// 添加文件输出 - 使用JSON格式
	if hasFile {
		fileOutput := &lumberjack.Logger{
			Filename:   config.FileConfig.Path,
			MaxSize:    config.FileConfig.MaxSize,
			MaxAge:     config.FileConfig.MaxAge,
			MaxBackups: config.FileConfig.MaxBackups,
			Compress:   config.FileConfig.Compress,
		}
		fileCore := createCore(fileOutput, level, JSONFormat, config)
		cores = append(cores, fileCore)
	}

	// 如果没有输出，默认输出到控制台
	if len(cores) == 0 {
		consoleCore := createCore(os.Stdout, level, TextFormat, config)
		cores = append(cores, consoleCore)
	}

	// 合并所有core
	var core zapcore.Core
	if len(cores) == 1 {
		core = cores[0]
	} else {
		core = zapcore.NewTee(cores...)
	}

	// 创建 Logger 选项
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	// 添加堆栈跟踪
	options = append(options, zap.AddStacktrace(stackLevel))

	// 开发模式设置
	if config.Development {
		options = append(options, zap.Development())
	}

	// 创建 Logger
	structLogger := zap.New(core, options...)

	return &ZapLogger{
		logger:        structLogger.Sugar(),
		structLogger:  structLogger,
		level:         level,
		format:        TextFormat, // 默认格式设为text
		style:         config.Style,
		development:   config.Development,
		sampling:      config.Sampling,
		encoderConfig: config.EncoderConfig,
	}
}

// createCore 创建一个zapcore.Core
func createCore(w io.Writer, level Level, format Format, config logconfig.LogConfig) zapcore.Core {
	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          config.EncoderConfig.TimeKey,
		LevelKey:         config.EncoderConfig.LevelKey,
		NameKey:          config.EncoderConfig.NameKey,
		CallerKey:        config.EncoderConfig.CallerKey,
		MessageKey:       config.EncoderConfig.MessageKey,
		StacktraceKey:    config.EncoderConfig.StacktraceKey,
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      getZapLevelEncoder(config.EncoderConfig.EncodeLevel),
		EncodeTime:       getZapTimeEncoder(config.EncoderConfig.EncodeTime),
		EncodeDuration:   getZapDurationEncoder(config.EncoderConfig.EncodeDuration),
		EncodeCaller:     getZapCallerEncoder(config.EncoderConfig.EncodeCaller),
		ConsoleSeparator: " ", // 添加控制台分隔符，让字段之间有空格
	}

	// 创建编码器
	var encoder zapcore.Encoder
	if format == JSONFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建 Core
	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(w),
		convertToZapLevel(level),
	)
}

// newZapLogger 创建一个新的 ZapLogger
func newZapLogger(w io.Writer, level Level, format Format, stackLevel *zapcore.Level) *ZapLogger {
	// 使用默认配置
	config := logconfig.DefaultLogConfig()
	config.Level = logconfig.LogLevel(level.String())
	if format == JSONFormat {
		config.Format = logconfig.LogFormatJSON
	} else {
		config.Format = logconfig.LogFormatText
	}

	// 创建core
	core := createCore(w, level, format, config)

	// 创建 Logger 选项
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	// 添加堆栈跟踪
	if stackLevel != nil {
		options = append(options, zap.AddStacktrace(*stackLevel))
	} else {
		// 默认在 Error 级别及以上添加堆栈跟踪
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// 开发模式设置
	if config.Development {
		options = append(options, zap.Development())
	}

	// 创建 Logger
	structLogger := zap.New(core, options...)

	return &ZapLogger{
		logger:        structLogger.Sugar(),
		structLogger:  structLogger,
		level:         level,
		format:        format,
		style:         config.Style,
		development:   config.Development,
		sampling:      config.Sampling,
		encoderConfig: config.EncoderConfig,
	}
}

// Debug 实现 Logger 接口
func (l *ZapLogger) Debug(msg string, keysAndValues ...any) {
	l.logger.Debugw(msg, keysAndValues...)
}

// Info 实现 Logger 接口
func (l *ZapLogger) Info(msg string, keysAndValues ...any) {
	l.logger.Infow(msg, keysAndValues...)
}

// Warn 实现 Logger 接口
func (l *ZapLogger) Warn(msg string, keysAndValues ...any) {
	l.logger.Warnw(msg, keysAndValues...)
}

// Error 实现 Logger 接口
func (l *ZapLogger) Error(msg string, keysAndValues ...any) {
	l.logger.Errorw(msg, keysAndValues...)
}

// Fatal 实现 Logger 接口
func (l *ZapLogger) Fatal(msg string, keysAndValues ...any) {
	l.logger.Fatalw(msg, keysAndValues...)
}

// WithFields 实现 Logger 接口
func (l *ZapLogger) WithFields(fields map[string]any) Logger {
	// 将 map 转换为 key-value 对
	keyValues := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		keyValues = append(keyValues, k, v)
	}

	// 创建新的 Logger
	newLogger := l.logger.With(keyValues...)

	// 创建 zap 字段
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return &ZapLogger{
		logger:        newLogger,
		structLogger:  l.structLogger.With(zapFields...),
		level:         l.level,
		format:        l.format,
		style:         l.style,
		development:   l.development,
		sampling:      l.sampling,
		encoderConfig: l.encoderConfig,
	}
}

// WithField 实现 Logger 接口
func (l *ZapLogger) WithField(key string, value any) Logger {
	newLogger := l.logger.With(key, value)

	return &ZapLogger{
		logger:        newLogger,
		structLogger:  l.structLogger.With(zap.Any(key, value)),
		level:         l.level,
		format:        l.format,
		style:         l.style,
		development:   l.development,
		sampling:      l.sampling,
		encoderConfig: l.encoderConfig,
	}
}

// SetLevel 实现 Logger 接口
func (l *ZapLogger) SetLevel(level Level) {
	l.level = level
	// 注意：zap 不支持动态修改日志级别，这里只是记录级别
}

// GetLevel 实现 Logger 接口
func (l *ZapLogger) GetLevel() Level {
	return l.level
}

// SetOutput 实现 Logger 接口
func (l *ZapLogger) SetOutput(w io.Writer) {
	// 注意：zap 不支持动态修改输出，这里是一个空操作
}

// SetFormat 实现 Logger 接口
func (l *ZapLogger) SetFormat(format Format) {
	l.format = format
	// 注意：zap 不支持动态修改格式，这里只是记录格式
}

// WithContext 实现 Logger 接口
func (l *ZapLogger) WithContext(ctx context.Context) Logger {
	// 从上下文中提取关键信息
	fields := extractContextFields(ctx)

	// 将字段添加到日志中
	keyValues := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		keyValues = append(keyValues, k, v)
	}

	// 创建新的 Logger
	newLogger := l.logger.With(keyValues...)

	// 创建 zap 字段
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return &ZapLogger{
		logger:        newLogger,
		structLogger:  l.structLogger.With(zapFields...),
		level:         l.level,
		format:        l.format,
		style:         l.style,
		development:   l.development,
		sampling:      l.sampling,
		encoderConfig: l.encoderConfig,
	}
}

// DebugContext 实现 Logger 接口
func (l *ZapLogger) DebugContext(ctx context.Context, msg string, keysAndValues ...any) {
	if l.style == logconfig.LogStyleStructured {
		// 使用结构化日志
		fields := extractContextFields(ctx)
		zapFields := make([]zap.Field, 0, len(fields)+len(keysAndValues)/2)

		// 添加上下文字段
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}

		// 添加用户提供的字段
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				key, ok := keysAndValues[i].(string)
				if !ok {
					key = "unknown"
				}
				zapFields = append(zapFields, zap.Any(key, keysAndValues[i+1]))
			}
		}

		l.structLogger.Debug(msg, zapFields...)
	} else {
		// 使用非结构化日志
		l.WithContext(ctx).Debug(msg, keysAndValues...)
	}
}

// InfoContext 实现 Logger 接口
func (l *ZapLogger) InfoContext(ctx context.Context, msg string, keysAndValues ...any) {
	if l.style == logconfig.LogStyleStructured {
		// 使用结构化日志
		fields := extractContextFields(ctx)
		zapFields := make([]zap.Field, 0, len(fields)+len(keysAndValues)/2)

		// 添加上下文字段
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}

		// 添加用户提供的字段
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				key, ok := keysAndValues[i].(string)
				if !ok {
					key = "unknown"
				}
				zapFields = append(zapFields, zap.Any(key, keysAndValues[i+1]))
			}
		}

		l.structLogger.Info(msg, zapFields...)
	} else {
		// 使用非结构化日志
		l.WithContext(ctx).Info(msg, keysAndValues...)
	}
}

// WarnContext 实现 Logger 接口
func (l *ZapLogger) WarnContext(ctx context.Context, msg string, keysAndValues ...any) {
	if l.style == logconfig.LogStyleStructured {
		// 使用结构化日志
		fields := extractContextFields(ctx)
		zapFields := make([]zap.Field, 0, len(fields)+len(keysAndValues)/2)

		// 添加上下文字段
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}

		// 添加用户提供的字段
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				key, ok := keysAndValues[i].(string)
				if !ok {
					key = "unknown"
				}
				zapFields = append(zapFields, zap.Any(key, keysAndValues[i+1]))
			}
		}

		l.structLogger.Warn(msg, zapFields...)
	} else {
		// 使用非结构化日志
		l.WithContext(ctx).Warn(msg, keysAndValues...)
	}
}

// ErrorContext 实现 Logger 接口
func (l *ZapLogger) ErrorContext(ctx context.Context, msg string, keysAndValues ...any) {
	if l.style == logconfig.LogStyleStructured {
		// 使用结构化日志
		fields := extractContextFields(ctx)
		zapFields := make([]zap.Field, 0, len(fields)+len(keysAndValues)/2)

		// 添加上下文字段
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}

		// 添加用户提供的字段
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				key, ok := keysAndValues[i].(string)
				if !ok {
					key = "unknown"
				}
				zapFields = append(zapFields, zap.Any(key, keysAndValues[i+1]))
			}
		}

		l.structLogger.Error(msg, zapFields...)
	} else {
		// 使用非结构化日志
		l.WithContext(ctx).Error(msg, keysAndValues...)
	}
}

// FatalContext 实现 Logger 接口
func (l *ZapLogger) FatalContext(ctx context.Context, msg string, keysAndValues ...any) {
	if l.style == logconfig.LogStyleStructured {
		// 使用结构化日志
		fields := extractContextFields(ctx)
		zapFields := make([]zap.Field, 0, len(fields)+len(keysAndValues)/2)

		// 添加上下文字段
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}

		// 添加用户提供的字段
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				key, ok := keysAndValues[i].(string)
				if !ok {
					key = "unknown"
				}
				zapFields = append(zapFields, zap.Any(key, keysAndValues[i+1]))
			}
		}

		l.structLogger.Fatal(msg, zapFields...)
	} else {
		// 使用非结构化日志
		l.WithContext(ctx).Fatal(msg, keysAndValues...)
	}
}

// extractContextFields 从上下文中提取字段
func extractContextFields(ctx context.Context) map[string]any {
	fields := make(map[string]any)

	// 提取请求ID
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		fields["request_id"] = requestID
	}

	// 提取用户ID
	if userID, ok := ctx.Value("user_id").(string); ok && userID != "" {
		fields["user_id"] = userID
	} else if userID, ok := ctx.Value("user_id").(int64); ok && userID != 0 {
		fields["user_id"] = userID
	} else if userID, ok := ctx.Value("user_id").(int); ok && userID != 0 {
		fields["user_id"] = userID
	}

	// 提取请求路径
	if path, ok := ctx.Value("path").(string); ok && path != "" {
		fields["path"] = path
	}

	// 提取请求方法
	if method, ok := ctx.Value("method").(string); ok && method != "" {
		fields["method"] = method
	}

	// 提取请求IP
	if ip, ok := ctx.Value("ip").(string); ok && ip != "" {
		fields["ip"] = ip
	}

	// 提取请求耗时
	if latency, ok := ctx.Value("latency").(time.Duration); ok {
		fields["latency_ms"] = latency.Milliseconds()
	}

	// 提取请求状态码
	if status, ok := ctx.Value("status").(int); ok {
		fields["status"] = status
	}

	// 提取请求用户代理
	if userAgent, ok := ctx.Value("user_agent").(string); ok && userAgent != "" {
		fields["user_agent"] = userAgent
	}

	// 提取请求来源
	if referer, ok := ctx.Value("referer").(string); ok && referer != "" {
		fields["referer"] = referer
	}

	// 提取请求体大小
	if bodySize, ok := ctx.Value("body_size").(int); ok {
		fields["body_size"] = bodySize
	}

	// 提取请求追踪信息
	if traceID, ok := ctx.Value("trace_id").(string); ok && traceID != "" {
		fields["trace_id"] = traceID
	}

	// 提取请求跨度信息
	if spanID, ok := ctx.Value("span_id").(string); ok && spanID != "" {
		fields["span_id"] = spanID
	}

	return fields
}

// convertToZapLevel 将我们的日志级别转换为 zap 的日志级别
func convertToZapLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// getZapLevelEncoder 获取级别编码器
func getZapLevelEncoder(encoderType string) zapcore.LevelEncoder {
	switch encoderType {
	case "capital":
		return zapcore.CapitalLevelEncoder
	case "color":
		return zapcore.CapitalColorLevelEncoder
	case "lowercase":
		return zapcore.LowercaseLevelEncoder
	default:
		return zapcore.CapitalLevelEncoder
	}
}

// getZapTimeEncoder 获取时间编码器
func getZapTimeEncoder(encoderType string) zapcore.TimeEncoder {
	switch encoderType {
	case "iso8601":
		return zapcore.ISO8601TimeEncoder
	case "rfc3339":
		return zapcore.RFC3339TimeEncoder
	case "rfc3339nano":
		return zapcore.RFC3339NanoTimeEncoder
	case "epoch":
		return zapcore.EpochTimeEncoder
	case "epochmilli":
		return zapcore.EpochMillisTimeEncoder
	case "epochnano":
		return zapcore.EpochNanosTimeEncoder
	default:
		return zapcore.ISO8601TimeEncoder
	}
}

// getZapDurationEncoder 获取时间间隔编码器
func getZapDurationEncoder(encoderType string) zapcore.DurationEncoder {
	switch encoderType {
	case "string":
		return zapcore.StringDurationEncoder
	case "nanos":
		return zapcore.NanosDurationEncoder
	case "ms":
		return func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendFloat64(float64(d) / float64(time.Millisecond))
		}
	default:
		return zapcore.StringDurationEncoder
	}
}

// getZapCallerEncoder 获取调用者编码器
func getZapCallerEncoder(encoderType string) zapcore.CallerEncoder {
	switch encoderType {
	case "short":
		return zapcore.ShortCallerEncoder
	case "full":
		return zapcore.FullCallerEncoder
	default:
		return zapcore.ShortCallerEncoder
	}
}
