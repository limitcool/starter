# ErrorX - 精简错误处理方案

ErrorX 提供了一个简洁而强大的错误处理方案，它能够追踪错误的完整传播路径，同时保持代码的简洁性和可读性。

## 设计理念

- **简洁性**：API 简单，使用方便，只有少量核心函数
- **完整链条**：记录错误的完整传播路径，包含位置信息
- **职责分离**：
  - Repository 和 Service 层只负责返回错误，不打印日志
  - Controller 层负责处理错误，打印日志，返回用户友好的错误信息
- **一次日志**：在 Controller 层只打印一次详细的错误日志，避免重复

## 核心组件

### 错误类型

- **AppError**: 应用程序错误的基本结构，包含错误码、错误消息和HTTP状态码

### 错误处理函数

- **Errorf(baseErr, format, args...)**: 创建一个格式化的应用程序错误
- **WrapError(err, message)**: 包装错误并添加上下文信息和位置信息
- **FormatErrorChain(err)**: 格式化错误链，包括位置信息
- **GetUserMessage(err)**: 获取用户友好的错误消息
- **GetErrorCode(err)**: 获取错误码
- **GetHttpStatus(err)**: 获取HTTP状态码
- **Is(err, target)**: 检查错误是否是特定的应用程序错误

### 响应函数

- **response.Success(ctx, data)**: 返回成功响应
- **response.Error(ctx, err)**: 处理错误，打印日志，返回用户友好的错误信息

## 使用方式

### Repository 层

Repository 层负责返回错误，不打印日志：

```go
// 获取用户示例
func (r *UserRepo) GetByID(id int64) (*model.User, error) {
    var user model.User
    err := r.DB.First(&user, id).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        // 创建错误并添加位置信息
        notFoundErr := errorx.Errorf(errorx.ErrNotFound, "用户ID %d 不存在", id)
        return nil, errorx.WrapError(notFoundErr, "")
    }

    if err != nil {
        // 包装错误并添加位置信息
        return nil, errorx.WrapError(err, "查询用户失败")
    }

    return &user, nil
}
```

### Service 层

Service 层负责处理业务逻辑，可能传递 Repository 层的错误，也可能产生自己的错误：

```go
// 用户资料示例
func (s *UserService) GetUserProfile(id int64) (*v1.UserProfile, error) {
    // 1. 调用 Repository 层获取用户
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        // 传递 Repository 层的错误，添加上下文信息
        return nil, errorx.WrapError(err, "获取用户资料失败")
    }

    // 2. Service 层的业务逻辑验证
    if user.Status == enum.UserStatusDisabled {
        // Service 层产生的错误
        disabledErr := errorx.Errorf(errorx.ErrUserDisabled, "用户账号已被禁用")
        return nil, errorx.WrapError(disabledErr, "")
    }

    // 3. 构建用户资料
    profile := &v1.UserProfile{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
    }

    return profile, nil
}
```

### Controller 层

Controller 层负责处理错误，打印日志，返回用户友好的错误信息：

```go
// 用户资料接口示例
func (c *UserController) GetProfile(ctx *gin.Context) {
    // 获取路径参数
    idStr := ctx.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        // 包装错误并添加位置信息
        err = errorx.WrapError(err, "无效的ID参数")
        // 处理错误，打印日志，返回用户友好的错误信息
        response.Error(ctx, err)
        return
    }

    // 调用服务层
    profile, err := c.userService.GetUserProfile(id)
    if err != nil {
        // 包装错误并添加位置信息
        err = errorx.WrapError(err, "处理获取用户资料请求失败")
        // 处理错误，打印日志，返回用户友好的错误信息
        response.Error(ctx, err)
        return
    }

    // 成功响应
    response.Success(ctx, profile)
}
```

## 错误链示例

当发生错误时，`response.Error` 函数会记录一次详细的错误日志，包含完整的错误链：

```log
ERROR 2023-05-15T10:15:23+08:00 API error occurred
  request_id: req-123456
  trace_id: trace-123456
  error_code: 404001
  error_chain: 处理获取用户资料请求失败 [internal/controller/user_controller.go:45]: 获取用户资料失败 [internal/service/user_service.go:67]: 查询用户失败 [internal/repository/user_repo.go:34]: 用户ID 123 不存在 [internal/repository/user_repo.go:30]
  path: /api/v1/users/123
  method: GET
  client_ip: 192.168.1.1
```

这样，我们可以清楚地看到错误的传播路径和每个错误发生的具体位置，大大提高了调试效率。

## 最佳实践

1. **在所有层**：
   - 使用 `WrapError` 包装错误并添加位置信息
   - 提供有意义的上下文消息，描述当前操作

2. **在 Repository 层**：
   - 对于数据库错误，使用 `WrapError` 添加具体的查询上下文
   - 对于"记录不存在"等常见错误，使用预定义的错误类型

3. **在 Service 层**：
   - 对于业务逻辑验证失败，创建新的错误并使用 `WrapError` 添加位置信息
   - 对于传递自 Repository 层或其他 Service 的错误，使用 `WrapError` 添加上下文

4. **在 Controller 层**：
   - 使用 `WrapError` 添加请求上下文
   - 使用 `response.Error` 处理错误，打印日志，返回用户友好的错误信息

5. **错误消息**：
   - 技术消息应该详细，包含关键参数，如 `用户ID 123 不存在`
   - 用户消息应该友好，不包含技术细节，如 `找不到该用户信息`

## 总结

这个精简的错误处理方案提供了一种简洁而强大的错误处理方式，它能够追踪错误的完整传播路径，同时保持代码的简洁性和可读性。通过职责分离和一次日志的原则，它避免了重复日志和代码混乱，是一种平衡了简洁性和功能性的错误处理方案。
