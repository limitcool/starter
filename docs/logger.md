# 日志系统使用指南

本项目使用了统一的日志抽象层，支持多种日志实现，目前支持 charmbracelet/log 和 uber-go/zap。

## 基本用法

```go
import "github.com/limitcool/starter/internal/pkg/logger"

// 记录不同级别的日志
logger.Debug("这是一条调试日志", "key1", "value1", "key2", 123)
logger.Info("这是一条信息日志", "user", "admin", "action", "login")
logger.Warn("这是一条警告日志", "latency", 200, "threshold", 100)
logger.Error("这是一条错误日志", "error", err, "request_id", "req-123")
logger.Fatal("这是一条致命错误日志") // 会导致程序退出

// 记录错误日志的辅助函数
logger.LogError("操作失败", err, "user_id", 123, "action", "update")
```

## 结构化日志

所有日志函数都支持键值对形式的结构化日志：

```go
logger.Info("用户登录",
    "user_id", 123,
    "username", "admin",
    "ip", "192.168.1.1",
    "duration_ms", 45,
)
```

## 创建带上下文的日志记录器

```go
// 创建带有字段的日志记录器
userLogger := logger.WithFields(map[string]interface{}{
    "user_id": 123,
    "username": "admin",
})

// 使用带上下文的日志记录器
userLogger.Info("用户登录成功")
userLogger.Error("用户操作失败", "error", err)

// 添加单个字段
requestLogger := logger.WithField("request_id", "req-123")
requestLogger.Info("开始处理请求")
```

## 错误处理与日志记录

结合 errorx 包使用：

```go
import (
    "github.com/limitcool/starter/internal/pkg/errorx"
    "github.com/limitcool/starter/internal/pkg/logger"
)

func doSomething() error {
    // 创建错误
    if err := someOperation(); err != nil {
        // 包装错误并添加上下文
        return errorx.WrapError(err, "操作失败")
    }
    return nil
}

// 在控制器中处理错误
func (c *Controller) HandleRequest(ctx *gin.Context) {
    if err := doSomething(); err != nil {
        // 记录错误，包含完整的错误链
        logger.LogErrorWithStack("处理请求失败", err, 
            "request_id", ctx.GetString("request_id"),
            "user_id", ctx.GetInt64("user_id"),
        )
        
        // 返回用户友好的错误响应
        response.Error(ctx, err)
        return
    }
    
    response.Success(ctx, "操作成功")
}
```

## 切换日志实现

项目默认使用 charmbracelet/log 作为日志实现，如果需要切换到 zap，可以在应用初始化时进行设置：

```go
import (
    "os"
    "github.com/limitcool/starter/internal/pkg/logger"
)

// 使用 Zap 作为日志实现
zapLogger := logger.NewZapLogger(os.Stdout, logger.InfoLevel, logger.JSONFormat)
logger.SetDefault(zapLogger)

// 现在所有的日志都会使用 Zap 记录
logger.Info("应用启动成功", "version", "1.0.0")
```

## 配置日志

在 `configs/config.yaml` 中配置日志：

```yaml
Log:
  Level: info                 # 日志级别: debug, info, warn, error
  Output: [console, file]     # 输出方式: console, file
  Format: text                # 日志格式: text, json
  FileConfig:
    Path: ./logs/app.log      # 日志文件路径
    MaxSize: 100              # 每个日志文件的最大大小(MB)
    MaxAge: 7                 # 保留日志文件的天数
    MaxBackups: 10            # 保留的旧日志文件的最大数量
    Compress: true            # 是否压缩旧的日志文件
  StackTraceEnabled: true     # 是否启用堆栈跟踪
  StackTraceLevel: error      # 记录堆栈的最低日志级别
  MaxStackFrames: 64          # 堆栈帧最大数量
```

## 最佳实践

1. **使用结构化日志**：始终使用键值对形式记录日志，而不是使用格式化字符串。
2. **合理选择日志级别**：
   - Debug：详细的调试信息，仅在开发环境使用
   - Info：常规操作信息，如请求处理、用户登录等
   - Warn：潜在问题或即将发生的问题，如性能下降、接近限制等
   - Error：错误信息，如请求处理失败、数据库错误等
   - Fatal：致命错误，会导致程序退出
3. **包含上下文信息**：日志中应包含足够的上下文信息，如用户ID、请求ID、操作类型等。
4. **避免敏感信息**：不要在日志中记录密码、令牌等敏感信息。
5. **使用 LogErrorWithStack**：记录错误时使用 `LogErrorWithStack` 函数，它会自动包含错误链和堆栈信息。
