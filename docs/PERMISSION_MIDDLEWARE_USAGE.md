# 权限中间件使用指南

## 🎯 设计理念

我们的权限中间件实现了声明式的权限控制，让API路由定义本身就能清晰地说明访问它需要什么权限。

## 🔧 中间件方法

### 1. RequirePermissionKey - 基于权限标识的控制

这是推荐的方式，直接使用权限字典中的Key：

```go
// 需要特定权限Key的中间件
func (m *PermissionMiddleware) RequirePermissionKey(permissionKey string) gin.HandlerFunc
```

**使用示例：**
```go
// 创建用户 - 需要 "user:create" 权限
v1.POST("/users", permissionMiddleware.RequirePermissionKey("user:create"), userHandler.Create)

// 查看用户列表 - 需要 "user:list" 权限  
v1.GET("/users", permissionMiddleware.RequirePermissionKey("user:list"), userHandler.List)

// 删除用户 - 需要 "user:delete" 权限
v1.DELETE("/users/:id", permissionMiddleware.RequirePermissionKey("user:delete"), userHandler.Delete)

// 查看课程列表 - 需要 "course:list" 权限
v1.GET("/courses", permissionMiddleware.RequirePermissionKey("course:list"), courseHandler.List)

// 创建课程 - 需要 "course:create" 权限
v1.POST("/courses", permissionMiddleware.RequirePermissionKey("course:create"), courseHandler.Create)
```

### 2. RequirePermission - 基于资源和操作的控制

传统方式，分别指定资源和操作：

```go
// 需要特定资源和操作权限的中间件
func (m *PermissionMiddleware) RequirePermission(resource, action string) gin.HandlerFunc
```

**使用示例：**
```go
// 等价于 RequirePermissionKey("user:create")
v1.POST("/users", permissionMiddleware.RequirePermission("user", "create"), userHandler.Create)

// 等价于 RequirePermissionKey("user:list")
v1.GET("/users", permissionMiddleware.RequirePermission("user", "list"), userHandler.List)
```

### 3. RequireAdmin - 管理员权限控制

```go
// 需要管理员权限的中间件
func (m *PermissionMiddleware) RequireAdmin() gin.HandlerFunc
```

**使用示例：**
```go
// 系统配置接口 - 只有管理员可以访问
v1.GET("/system/config", permissionMiddleware.RequireAdmin(), systemHandler.GetConfig)
v1.PUT("/system/config", permissionMiddleware.RequireAdmin(), systemHandler.UpdateConfig)
```

### 4. RequireRole - 基于角色的控制

```go
// 需要特定角色的中间件
func (m *PermissionMiddleware) RequireRole(roleName string) gin.HandlerFunc
```

**使用示例：**
```go
// 教练专用接口
v1.GET("/coach/students", permissionMiddleware.RequireRole("教练"), coachHandler.GetStudents)

// 销售专用接口
v1.GET("/sales/leads", permissionMiddleware.RequireRole("销售"), salesHandler.GetLeads)
```

## 🚀 完整路由示例

```go
package router

import (
    "github.com/gin-gonic/gin"
    "github.com/limitcool/starter/internal/handler"
    "github.com/limitcool/starter/internal/middleware"
    "github.com/limitcool/starter/internal/pkg/permission"
)

func SetupRoutes(
    router *gin.Engine,
    permissionService *permission.Service,
    userHandler *handler.UserHandler,
    courseHandler *handler.CourseHandler,
    memberHandler *handler.MemberHandler,
) {
    // 创建权限中间件
    permissionMiddleware := middleware.NewPermissionMiddleware(permissionService)
    
    // API v1 路由组
    v1 := router.Group("/api/v1")
    
    // 需要认证的路由组
    auth := v1.Group("")
    auth.Use(middleware.JWTMiddleware()) // JWT认证中间件
    
    // 用户管理 API
    userGroup := auth.Group("/users")
    {
        // 查看用户列表 - 需要 user:list 权限
        userGroup.GET("", permissionMiddleware.RequirePermissionKey("user:list"), userHandler.List)
        
        // 创建用户 - 需要 user:create 权限
        userGroup.POST("", permissionMiddleware.RequirePermissionKey("user:create"), userHandler.Create)
        
        // 查看用户详情 - 需要 user:view 权限
        userGroup.GET("/:id", permissionMiddleware.RequirePermissionKey("user:view"), userHandler.Get)
        
        // 更新用户 - 需要 user:update 权限
        userGroup.PUT("/:id", permissionMiddleware.RequirePermissionKey("user:update"), userHandler.Update)
        
        // 删除用户 - 需要 user:delete 权限
        userGroup.DELETE("/:id", permissionMiddleware.RequirePermissionKey("user:delete"), userHandler.Delete)
    }
    
    // 课程管理 API
    courseGroup := auth.Group("/courses")
    {
        // 查看课程列表 - 需要 course:list 权限
        courseGroup.GET("", permissionMiddleware.RequirePermissionKey("course:list"), courseHandler.List)
        
        // 创建课程 - 需要 course:create 权限
        courseGroup.POST("", permissionMiddleware.RequirePermissionKey("course:create"), courseHandler.Create)
        
        // 更新课程 - 需要 course:update 权限
        courseGroup.PUT("/:id", permissionMiddleware.RequirePermissionKey("course:update"), courseHandler.Update)
        
        // 删除课程 - 需要 course:delete 权限
        courseGroup.DELETE("/:id", permissionMiddleware.RequirePermissionKey("course:delete"), courseHandler.Delete)
    }
    
    // 会员管理 API
    memberGroup := auth.Group("/members")
    {
        // 查看会员列表 - 需要 member:list 权限
        memberGroup.GET("", permissionMiddleware.RequirePermissionKey("member:list"), memberHandler.List)
        
        // 编辑会员信息 - 需要 member:edit 权限
        memberGroup.PUT("/:id", permissionMiddleware.RequirePermissionKey("member:edit"), memberHandler.Update)
    }
    
    // 管理员专用 API
    adminGroup := auth.Group("/admin")
    {
        // 系统配置 - 只有管理员可以访问
        adminGroup.GET("/config", permissionMiddleware.RequireAdmin(), systemHandler.GetConfig)
        adminGroup.PUT("/config", permissionMiddleware.RequireAdmin(), systemHandler.UpdateConfig)
        
        // 权限管理 - 需要系统管理权限
        adminGroup.GET("/permissions", permissionMiddleware.RequirePermissionKey("sys:permission"), permissionHandler.List)
        adminGroup.POST("/permissions/assign", permissionMiddleware.RequirePermissionKey("sys:permission"), permissionHandler.Assign)
    }
    
    // 小程序端 API
    mpGroup := auth.Group("/mp")
    {
        // 我的学员 - 需要 mp_student:list 权限（教练小程序端）
        mpGroup.GET("/students", permissionMiddleware.RequirePermissionKey("mp_student:list"), mpHandler.GetStudents)
    }
}
```

## 🔍 工作流程

### 1. 权限检查流程

```
请求API → JWT认证中间件 → 权限中间件 → 业务处理器
           ↓                ↓
       设置user_id      检查权限Key → Casbin验证 → 通过/拒绝
```

### 2. 权限验证逻辑

1. **获取用户ID**：从JWT认证中间件设置的`user_id`获取
2. **管理员检查**：如果是管理员，直接通过
3. **解析权限Key**：将`"user:create"`解析为`resource="user", action="create"`
4. **Casbin验证**：调用`CheckPermission(userID, resource, action)`
5. **返回结果**：通过则继续，否则返回403错误

### 3. 错误处理

- **401 Unauthorized**：用户未登录或token无效
- **403 Forbidden**：用户已登录但权限不足
- **500 Internal Server Error**：权限检查过程中发生错误

## 💡 最佳实践

### 1. 权限Key命名规范

- 格式：`resource:action`
- 资源名：使用单数形式，如`user`、`course`、`member`
- 操作名：使用动词，如`create`、`list`、`update`、`delete`、`view`

### 2. 路由组织

- 按功能模块分组
- 相同权限要求的接口放在一起
- 使用中间件链式调用

### 3. 权限粒度

- **粗粒度**：适用于简单场景，如`user:manage`
- **细粒度**：适用于复杂场景，如`user:create`、`user:update`、`user:delete`

### 4. 特殊权限

- **管理员权限**：使用`RequireAdmin()`
- **角色权限**：使用`RequireRole(roleName)`
- **平台权限**：使用不同的权限Key区分，如`mp_student:list`

## 🎉 总结

通过这种声明式的权限控制方式，我们实现了：

1. **高可读性**：路由定义即权限文档
2. **易维护性**：权限变更只需修改中间件参数
3. **解耦设计**：业务逻辑无需关心权限检查
4. **灵活配置**：支持多种权限控制方式
5. **统一管理**：所有权限检查逻辑集中在中间件中

这样的设计让权限控制变得简单、清晰、易于维护！🚀
