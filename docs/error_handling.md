# 统一错误处理文档

## 概述

本文档描述了一种统一的错误处理模式，其中服务层和模型层只负责上抛错误，而在中间件中统一处理错误、记录日志和返回响应。这种模式有助于保持代码的一致性和可维护性。

## 错误处理原则

1. **单一职责原则**：
   - 仓库层负责数据访问和基本错误包装
   - 服务层负责业务逻辑和错误上抛
   - 控制器层负责参数验证和错误上抛
   - 中间件层负责统一处理错误、记录日志和返回响应

2. **错误包装原则**：
   - 使用 `errorx.WrapError` 包装错误，添加上下文信息
   - 保留原始错误，形成错误链
   - 添加位置信息，便于调试

3. **错误记录原则**：
   - 只在中间件中记录错误日志，避免重复记录
   - 记录完整的错误链和位置信息
   - 记录请求上下文信息，如请求ID、用户ID等

4. **错误响应原则**：
   - 向用户返回友好的错误信息
   - 隐藏敏感的技术细节
   - 保持响应格式的一致性

## 错误处理中间件

### 中间件实现

```go
// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 使用defer+recover捕获所有可能的panic
        defer func() {
            if err := recover(); err != nil {
                // 记录堆栈信息
                stack := string(debug.Stack())

                // 获取请求上下文
                ctx := c.Request.Context()

                // 记录详细的panic日志
                logger.ErrorContext(ctx, "Panic recovered",
                    "error", err,
                    "stack", stack,
                    "path", c.Request.URL.Path,
                    "method", c.Request.Method,
                    "client_ip", c.ClientIP())

                // 根据不同类型的panic返回不同的错误
                var appErr *errorx.AppError
                switch e := err.(type) {
                case *errorx.AppError:
                    appErr = e
                case error:
                    appErr = errorx.ErrInternal.WithError(e)
                case string:
                    appErr = errorx.ErrInternal.WithMsg(e)
                default:
                    appErr = errorx.ErrInternal.WithMsg(fmt.Sprintf("%v", err))
                }

                // 返回错误响应
                response.Error(c, appErr)
                c.Abort()
            }
        }()

        // 处理请求
        c.Next()

        // 检查是否有错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            handleError(c, err)
            c.Abort()
        }
    }
}

// handleError 处理不同类型的错误
func handleError(c *gin.Context, err error) {
    // 获取请求上下文
    ctx := c.Request.Context()

    // 使用统一的错误响应函数
    // response.Error 内部会记录错误日志，所以这里不需要重复记录
    response.Error(c, err)
}
```

### 错误响应函数

```go
// Error 返回错误响应
func Error(c *gin.Context, err error, msg ...string) {
    var (
        httpStatus int
        errorCode  int
        message    string
    )

    // 获取请求上下文
    ctx := c.Request.Context()

    // 尝试使用错误链推导错误码
    err = errorx.WrapErrorWithContext(ctx, err, "")

    // 获取错误信息
    if appErr, ok := err.(*errorx.AppError); ok {
        // 如果是 AppError类型，直接使用其属性
        message = appErr.GetErrorMsg()
        httpStatus = getHttpStatus(appErr)
        errorCode = appErr.GetErrorCode()
    } else {
        // 如果不是 AppError类型，使用默认值
        message = err.Error()
        httpStatus = http.StatusInternalServerError
        errorCode = errorx.ErrorUnknownCode
    }

    // 允许调用方覆盖原始错误消息
    if len(msg) > 0 {
        message = msg[0]
    }

    // 记录错误到日志
    logger.ErrorContext(ctx, "API error occurred",
        "code", errorCode,
        "message", message,
        "path", c.Request.URL.Path,
        "method", c.Request.Method,
        "client_ip", c.ClientIP(),
        "error_chain", errorx.FormatErrorChainWithContext(ctx, err),
    )

    // 统一响应结构
    c.JSON(httpStatus, Response{
        Code:    errorCode,
        Msg:     message,
        Data:    struct{}{},
        ReqID:   getRequestID(c),
        Time:    time.Now().Unix(),
        TraceID: getTraceIDFromContext(c),
    })
}
```

## 错误包装函数

```go
// WrapError 包装错误并添加上下文信息和位置信息
// 用法: WrapError(err, "查询用户失败")
// 如果不需要添加消息，可以传入空字符串: WrapError(err, "")
func WrapError(err error, message string) error {
    return WrapErrorWithContext(context.Background(), err, message)
}

// WrapErrorWithContext 使用上下文包装错误并添加上下文信息和位置信息
// 用法: WrapErrorWithContext(ctx, err, "查询用户失败")
// 如果不需要添加消息，可以传入空字符串: WrapErrorWithContext(ctx, err, "")
func WrapErrorWithContext(ctx context.Context, err error, message string) error {
    if err == nil {
        return nil
    }

    // 获取调用者的文件和行号
    _, file, line, ok := runtime.Caller(1)
    if !ok {
        file = "unknown"
        line = 0
    }

    // 简化文件路径，只保留最后几个部分
    parts := strings.Split(file, "/")
    if len(parts) > 3 {
        file = strings.Join(parts[len(parts)-3:], "/")
    }

    // 位置信息
    location := fmt.Sprintf("[%s:%d]", file, line)

    // 如果没有提供消息，只添加位置信息
    if message == "" {
        return errors.WithMessage(err, location)
    }

    // 添加位置信息和消息
    return errors.WithMessage(err, fmt.Sprintf("%s %s", message, location))
}
```

## 使用示例

### 控制器层

```go
func (c *UserController) GetUser(ctx *gin.Context) {
    // 获取用户ID
    idStr := ctx.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        // 参数验证错误，直接上抛
        ctx.Error(errorx.ErrInvalidParams.WithError(err))
        return
    }

    // 调用服务层方法
    user, err := c.userService.GetUserByID(ctx.Request.Context(), id)
    if err != nil {
        // 直接上抛错误，由错误处理中间件处理
        ctx.Error(err)
        return
    }

    // 成功响应
    response.Success(ctx, user)
}
```

### 服务层

```go
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
    // 调用仓库层方法
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        // 包装错误并上抛，添加上下文信息
        return nil, errorx.WrapError(err, fmt.Sprintf("获取用户失败: id=%d", id))
    }

    // 业务逻辑处理
    if !user.Enabled {
        return nil, errorx.ErrUserDisabled.WithMsg(fmt.Sprintf("用户 %s 已被禁用", user.Username))
    }

    return user, nil
}
```

### 仓库层

```go
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
    var user model.User
    err := r.DB.WithContext(ctx).First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 记录不存在，返回特定错误
            return nil, errorx.WrapError(
                errorx.NewAppError(r.ErrorCode, fmt.Sprintf("用户ID %d 不存在", id), http.StatusNotFound),
                "",
            )
        }
        // 其他数据库错误
        return nil, errorx.WrapError(err, fmt.Sprintf("查询用户失败: id=%d", id))
    }

    return &user, nil
}
```

## 错误类型

### 应用错误

```go
// AppError 是应用程序错误的基本结构
type AppError struct {
    code     int    // 错误码
    msg      string // 错误消息
    httpCode int    // HTTP状态码
    traceID  string // 链路追踪 ID
    orig     error  // 原始错误
}

// Error 实现 error 接口
func (e *AppError) Error() string {
    return e.msg
}

// GetErrorCode 获取错误码
func (e *AppError) GetErrorCode() int {
    return e.code
}

// GetErrorMsg 获取错误消息
func (e *AppError) GetErrorMsg() string {
    return e.msg
}

// GetHttpStatus 获取HTTP状态码
func (e *AppError) GetHttpStatus() int {
    return e.httpCode
}

// WithMsg 添加错误消息
func (e *AppError) WithMsg(msg string) *AppError {
    e.msg = msg
    return e
}

// WithError 添加原始错误
func (e *AppError) WithError(err error) *AppError {
    e.orig = err
    return e
}

// Unwrap 获取原始错误
func (e *AppError) Unwrap() error {
    return e.orig
}
```

### 预定义错误

```go
// 预定义错误实例
var (
    // 基础错误
    ErrSuccess = NewAppError(SuccessCode, "成功", http.StatusOK)

    // 通用错误
    ErrUnknown = NewAppError(ErrorUnknownCode, "未知错误", http.StatusInternalServerError)
    ErrInternal = NewAppError(ErrorInternalCode, "服务器内部错误", http.StatusInternalServerError)
    ErrInvalidParams = NewAppError(ErrorInvalidParamsCode, "无效的参数", http.StatusBadRequest)
    ErrNotFound = NewAppError(ErrorNotFoundCode, "资源不存在", http.StatusNotFound)
    ErrForbidden = NewAppError(ErrorForbiddenCode, "禁止访问", http.StatusForbidden)
    ErrUnauthorized = NewAppError(ErrorUnauthorizedCode, "未授权", http.StatusUnauthorized)
    ErrTimeout = NewAppError(ErrorTimeoutCode, "请求超时", http.StatusRequestTimeout)
    ErrTooManyRequests = NewAppError(ErrorTooManyRequestsCode, "请求过于频繁", http.StatusTooManyRequests)

    // 用户错误
    ErrUserNotFound = NewAppError(ErrorUserNotFoundCode, "用户不存在", http.StatusNotFound)
    ErrUserExists = NewAppError(ErrorUserExistsCode, "用户已存在", http.StatusConflict)
    ErrUserPasswordError = NewAppError(ErrorUserPasswordErrorCode, "用户密码错误", http.StatusUnauthorized)
    ErrUserNoLogin = NewAppError(ErrorUserNoLoginCode, "用户未登录", http.StatusUnauthorized)
    ErrUserDisabled = NewAppError(ErrorUserDisabledCode, "用户已禁用", http.StatusForbidden)
)
```

## 最佳实践

1. **仓库层**：
   - 使用 `errorx.WrapError` 包装数据库错误
   - 对特定错误（如记录不存在）返回特定的错误类型
   - 添加详细的上下文信息，如查询参数

2. **服务层**：
   - 使用 `errorx.WrapError` 包装仓库层错误
   - 添加业务相关的上下文信息
   - 对业务规则验证失败返回特定的错误类型

3. **控制器层**：
   - 使用 `ctx.Error` 上抛错误，由中间件统一处理
   - 对参数验证错误返回 `errorx.ErrInvalidParams`
   - 不要在控制器层记录错误日志，由中间件统一记录

4. **中间件层**：
   - 统一处理所有错误
   - 记录详细的错误日志
   - 返回用户友好的错误响应

## 总结

这种错误处理模式的优点是：

1. **统一处理**：所有错误都在中间件中统一处理，保持一致性
2. **清晰分离**：服务层和模型层只负责上抛错误，不需要关心错误如何被记录和返回
3. **详细日志**：在中间件中记录详细的错误日志，包括错误链和位置信息
4. **用户友好**：向用户返回友好的错误信息，同时在日志中保留详细的技术信息

通过这种模式，可以大大简化错误处理逻辑，提高代码的可维护性和一致性。
