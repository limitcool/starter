# 统一错误处理中间件

本文档介绍如何使用统一错误处理中间件来简化控制器代码，避免在每个控制器方法中重复调用 `response.Error()`。

## 背景

在传统的错误处理方式中，控制器方法需要在每个可能出错的地方都调用 `response.Error()`，这导致了大量的重复代码。例如：

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    id := cast.ToUint(c.Param("id"))

    user, err := h.userRepo.GetByID(c.Request.Context(), id)
    if err != nil {
        logger.Error("Failed to get user", "error", err, "id", id)
        response.Error(c, err)
        return
    }

    response.Success(c, user)
}
```

## 改进的统一错误处理方式

改进后的错误处理方式使用全局中间件捕获控制器方法返回的错误，而不是要求每个控制器方法都手动处理错误。这样，我们只需要在应用程序启动时注册一次中间件，而不是在每个控制器方法中重复处理错误。

### 1. 错误处理中间件

我们使用两个中间件来处理不同类型的错误：

1. **PanicRecovery**：用于捕获 panic 并返回友好的错误响应
2. **GlobalErrorHandler**：用于处理控制器方法通过 `c.Error()` 返回的错误

```go
// 在 router/module.go 中注册中间件
r.Use(middleware.PanicRecovery())
r.Use(middleware.GlobalErrorHandler())
```

### 2. 辅助函数 ErrorHandlerFunc

我们提供了一个辅助函数 `ErrorHandlerFunc`，用于将返回 `error` 的控制器方法转换为 `gin.HandlerFunc`：

```go
// ErrorHandlerFunc 是一个辅助函数，用于将返回 error 的控制器方法转换为 gin.HandlerFunc
func ErrorHandlerFunc(handler func(c *gin.Context) error) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := handler(c); err != nil {
            // 如果已经有响应写入，不再处理错误
            if c.Writer.Written() {
                return
            }
            
            // 将错误添加到 gin 的错误链中，由 GlobalErrorHandler 统一处理
            _ = c.Error(err)
        }
    }
}
```

### 3. 控制器方法改造

控制器方法需要返回 `error`，而不是直接处理错误：

```go
func (h *UserHandler) GetUser(c *gin.Context) error {
    id := cast.ToUint(c.Param("id"))

    user, err := h.userRepo.GetByID(c.Request.Context(), id)
    if err != nil {
        return err
    }

    response.Success(c, user)
    return nil
}
```

### 4. 路由注册

在路由注册时，使用 `ErrorHandlerFunc` 包装控制器方法：

```go
// 使用 ErrorHandlerFunc 包装控制器方法
r.GET("/users/:id", middleware.ErrorHandlerFunc(userHandler.GetUser))
r.POST("/users", middleware.ErrorHandlerFunc(userHandler.CreateUser))

// 对于需要中间件的路由，可以正常添加中间件
admin := r.Group("/admin", middleware.JWTAuth(config), middleware.AdminCheck())
{
    admin.GET("/users/:id", middleware.ErrorHandlerFunc(userHandler.GetUser))
    admin.POST("/users", middleware.ErrorHandlerFunc(userHandler.CreateUser))
}
```

## 优势

1. **全局处理**：错误处理逻辑集中在一处，便于修改和扩展
2. **代码更简洁**：控制器方法只需返回错误，不需要直接处理错误
3. **一致性**：所有错误都以相同的方式处理，确保响应格式一致
4. **更少的包装**：不需要为每个路由都包装处理函数，只需在路由注册时使用 `ErrorHandlerFunc`

## 注意事项

1. 控制器方法必须返回 `error` 类型
2. 成功响应后必须返回 `nil`
3. 不要在控制器方法中调用 `response.Error()`，而是返回错误
4. 如果需要自定义错误消息，使用 `errorx.WrapError()` 或 `errorx.Errorf()`
