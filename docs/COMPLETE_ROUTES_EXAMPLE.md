# 完整的权限系统路由配置示例

## 🎯 概述

这是一个完整的路由配置示例，展示了如何将权限管理、菜单管理和用户菜单接口整合在一起。

## 🔧 完整路由配置

```go
package router

import (
    "github.com/gin-gonic/gin"
    "github.com/limitcool/starter/internal/handler"
    "github.com/limitcool/starter/internal/middleware"
    "github.com/limitcool/starter/internal/model"
    "github.com/limitcool/starter/internal/pkg/permission"
)

// SetupPermissionRoutes 设置完整的权限系统路由
func SetupPermissionRoutes(
    router *gin.Engine,
    permissionService *permission.Service,
    permissionRepo *model.PermissionRepo,
    menuRepo *model.MenuRepo,
    roleRepo *model.RoleRepo,
) {
    // 创建权限中间件
    permissionMiddleware := middleware.NewPermissionMiddleware(permissionService)
    
    // 创建处理器
    permissionHandler := handler.NewPermissionHandler(permissionService)
    permissionCRUDHandler := handler.NewPermissionCRUDHandler(permissionRepo)
    menuCRUDHandler := handler.NewMenuCRUDHandler(menuRepo)
    userMenuHandler := handler.NewUserMenuHandler(permissionService)

    // API v1 路由组
    v1 := router.Group("/api/v1")
    
    // 需要认证的路由组
    auth := v1.Group("")
    auth.Use(middleware.JWTMiddleware()) // JWT认证中间件

    // ==================== 用户个人信息相关 API ====================
    userGroup := auth.Group("/user")
    {
        // 获取我的菜单 - 只需要登录即可
        userGroup.GET("/menus", userMenuHandler.GetMyMenus)

        // 获取我的权限 - 只需要登录即可
        userGroup.GET("/permissions", userMenuHandler.GetMyPermissions)

        // 获取我的角色 - 只需要登录即可
        userGroup.GET("/roles", userMenuHandler.GetMyRoles)

        // 检查我的权限 - 只需要登录即可
        userGroup.POST("/check-permission", userMenuHandler.CheckMyPermission)
    }

    // ==================== 权限管理相关 API ====================
    permissionGroup := auth.Group("/admin/permissions")
    {
        // 权限分配接口
        // 分配用户角色 - 需要 "sys:user" 权限
        permissionGroup.POST("/assign-user-roles", 
            permissionMiddleware.RequirePermissionKey("sys:user"), 
            permissionHandler.AssignUserRoles)

        // 分配角色权限 - 需要 "sys:permission" 权限
        permissionGroup.POST("/roles/assign-permissions", 
            permissionMiddleware.RequirePermissionKey("sys:permission"), 
            permissionHandler.AssignRolePermissions)

        // 获取角色列表 - 需要 "role:list" 权限
        permissionGroup.GET("/roles", 
            permissionMiddleware.RequirePermissionKey("role:list"), 
            permissionHandler.GetRoles)

        // 创建角色 - 需要 "role:create" 权限
        permissionGroup.POST("/roles", 
            permissionMiddleware.RequirePermissionKey("role:create"), 
            permissionHandler.CreateRole)

        // 删除角色 - 需要 "role:delete" 权限
        permissionGroup.DELETE("/roles/:id", 
            permissionMiddleware.RequirePermissionKey("role:delete"), 
            permissionHandler.DeleteRole)

        // ==================== 权限CRUD接口 ====================
        // 权限列表 - 需要 "permission:list" 权限
        permissionGroup.GET("", 
            permissionMiddleware.RequirePermissionKey("permission:list"), 
            permissionCRUDHandler.GetPermissions)

        // 权限树 - 需要 "permission:list" 权限
        permissionGroup.GET("/tree", 
            permissionMiddleware.RequirePermissionKey("permission:list"), 
            permissionCRUDHandler.GetPermissionTree)

        // 权限详情 - 需要 "permission:view" 权限
        permissionGroup.GET("/:id", 
            permissionMiddleware.RequirePermissionKey("permission:view"), 
            permissionCRUDHandler.GetPermission)

        // 创建权限 - 需要 "permission:create" 权限
        permissionGroup.POST("", 
            permissionMiddleware.RequirePermissionKey("permission:create"), 
            permissionCRUDHandler.CreatePermission)

        // 更新权限 - 需要 "permission:update" 权限
        permissionGroup.PUT("/:id", 
            permissionMiddleware.RequirePermissionKey("permission:update"), 
            permissionCRUDHandler.UpdatePermission)

        // 删除权限 - 需要 "permission:delete" 权限
        permissionGroup.DELETE("/:id", 
            permissionMiddleware.RequirePermissionKey("permission:delete"), 
            permissionCRUDHandler.DeletePermission)
    }

    // ==================== 菜单管理相关 API ====================
    menuGroup := auth.Group("/admin/menus")
    {
        // 菜单列表 - 需要 "menu:list" 权限
        menuGroup.GET("", 
            permissionMiddleware.RequirePermissionKey("menu:list"), 
            menuCRUDHandler.GetMenus)

        // 菜单树 - 需要 "menu:list" 权限
        menuGroup.GET("/tree", 
            permissionMiddleware.RequirePermissionKey("menu:list"), 
            menuCRUDHandler.GetMenuTree)

        // 菜单详情 - 需要 "menu:view" 权限
        menuGroup.GET("/:id", 
            permissionMiddleware.RequirePermissionKey("menu:view"), 
            menuCRUDHandler.GetMenu)

        // 创建菜单 - 需要 "menu:create" 权限
        menuGroup.POST("", 
            permissionMiddleware.RequirePermissionKey("menu:create"), 
            menuCRUDHandler.CreateMenu)

        // 更新菜单 - 需要 "menu:update" 权限
        menuGroup.PUT("/:id", 
            permissionMiddleware.RequirePermissionKey("menu:update"), 
            menuCRUDHandler.UpdateMenu)

        // 删除菜单 - 需要 "menu:delete" 权限
        menuGroup.DELETE("/:id", 
            permissionMiddleware.RequirePermissionKey("menu:delete"), 
            menuCRUDHandler.DeleteMenu)

        // 更新菜单排序 - 需要 "menu:update" 权限
        menuGroup.PUT("/sort", 
            permissionMiddleware.RequirePermissionKey("menu:update"), 
            menuCRUDHandler.UpdateMenuSort)
    }

    // ==================== 业务模块示例 ====================
    // 用户管理 API
    userManagementGroup := auth.Group("/users")
    {
        // 查看用户列表 - 需要 "user:list" 权限
        userManagementGroup.GET("", 
            permissionMiddleware.RequirePermissionKey("user:list"), 
            userHandler.List)

        // 创建用户 - 需要 "user:create" 权限
        userManagementGroup.POST("", 
            permissionMiddleware.RequirePermissionKey("user:create"), 
            userHandler.Create)

        // 查看用户详情 - 需要 "user:view" 权限
        userManagementGroup.GET("/:id", 
            permissionMiddleware.RequirePermissionKey("user:view"), 
            userHandler.GetByID)

        // 更新用户 - 需要 "user:update" 权限
        userManagementGroup.PUT("/:id", 
            permissionMiddleware.RequirePermissionKey("user:update"), 
            userHandler.Update)

        // 删除用户 - 需要 "user:delete" 权限
        userManagementGroup.DELETE("/:id", 
            permissionMiddleware.RequirePermissionKey("user:delete"), 
            userHandler.Delete)
    }

    // 课程管理 API
    courseGroup := auth.Group("/courses")
    {
        // 查看课程列表 - 需要 "course:list" 权限
        courseGroup.GET("", 
            permissionMiddleware.RequirePermissionKey("course:list"), 
            courseHandler.List)

        // 创建课程 - 需要 "course:create" 权限
        courseGroup.POST("", 
            permissionMiddleware.RequirePermissionKey("course:create"), 
            courseHandler.Create)

        // 更新课程 - 需要 "course:update" 权限
        courseGroup.PUT("/:id", 
            permissionMiddleware.RequirePermissionKey("course:update"), 
            courseHandler.Update)

        // 删除课程 - 需要 "course:delete" 权限
        courseGroup.DELETE("/:id", 
            permissionMiddleware.RequirePermissionKey("course:delete"), 
            courseHandler.Delete)
    }

    // 会员管理 API
    memberGroup := auth.Group("/members")
    {
        // 查看会员列表 - 需要 "member:list" 权限
        memberGroup.GET("", 
            permissionMiddleware.RequirePermissionKey("member:list"), 
            memberHandler.List)

        // 编辑会员信息 - 需要 "member:edit" 权限
        memberGroup.PUT("/:id", 
            permissionMiddleware.RequirePermissionKey("member:edit"), 
            memberHandler.Update)

        // 查看会员详情 - 需要 "member:view" 权限
        memberGroup.GET("/:id", 
            permissionMiddleware.RequirePermissionKey("member:view"), 
            memberHandler.GetByID)
    }

    // 小程序端 API
    mpGroup := auth.Group("/mp")
    {
        // 我的学员 - 需要 "mp_student:list" 权限（教练小程序端）
        mpGroup.GET("/students", 
            permissionMiddleware.RequirePermissionKey("mp_student:list"), 
            mpHandler.GetStudents)

        // 学员详情 - 需要 "mp_student:view" 权限
        mpGroup.GET("/students/:id", 
            permissionMiddleware.RequirePermissionKey("mp_student:view"), 
            mpHandler.GetStudentDetail)
    }

    // ==================== 管理员专用 API ====================
    adminGroup := auth.Group("/admin/system")
    {
        // 系统配置 - 只有管理员可以访问
        adminGroup.GET("/config", 
            permissionMiddleware.RequireAdmin(), 
            systemHandler.GetConfig)

        adminGroup.PUT("/config", 
            permissionMiddleware.RequireAdmin(), 
            systemHandler.UpdateConfig)

        // 系统日志 - 只有管理员可以访问
        adminGroup.GET("/logs", 
            permissionMiddleware.RequireAdmin(), 
            systemHandler.GetLogs)
    }

    // ==================== 角色特定的路由 ====================
    // 教练专用接口
    coachGroup := auth.Group("/coach")
    {
        // 教练仪表盘 - 需要"教练"角色
        coachGroup.GET("/dashboard", 
            permissionMiddleware.RequireRole("教练"), 
            coachHandler.GetDashboard)

        // 教练课程表 - 需要"教练"角色
        coachGroup.GET("/schedule", 
            permissionMiddleware.RequireRole("教练"), 
            coachHandler.GetSchedule)
    }

    // 销售专用接口
    salesGroup := auth.Group("/sales")
    {
        // 销售线索 - 需要"销售"角色
        salesGroup.GET("/leads", 
            permissionMiddleware.RequireRole("销售"), 
            salesHandler.GetLeads)

        // 创建线索 - 需要"销售"角色
        salesGroup.POST("/leads", 
            permissionMiddleware.RequireRole("销售"), 
            salesHandler.CreateLead)
    }

    // ==================== 公开API ====================
    publicGroup := v1.Group("/public")
    {
        // 健康检查 - 无需权限
        publicGroup.GET("/health", healthHandler.Check)
        
        // 版本信息 - 无需权限
        publicGroup.GET("/version", versionHandler.Get)
    }
}
```

## 📊 权限Key与API的完整映射

### 系统管理权限
- `sys:user` → 用户角色分配
- `sys:permission` → 角色权限分配

### 权限管理权限
- `permission:list` → 查看权限列表/树
- `permission:view` → 查看权限详情
- `permission:create` → 创建权限
- `permission:update` → 更新权限
- `permission:delete` → 删除权限

### 菜单管理权限
- `menu:list` → 查看菜单列表/树
- `menu:view` → 查看菜单详情
- `menu:create` → 创建菜单
- `menu:update` → 更新菜单/排序
- `menu:delete` → 删除菜单

### 角色管理权限
- `role:list` → 查看角色列表
- `role:create` → 创建角色
- `role:delete` → 删除角色

### 业务模块权限
- `user:list/view/create/update/delete` → 用户管理
- `course:list/create/update/delete` → 课程管理
- `member:list/view/edit` → 会员管理
- `mp_student:list/view` → 小程序学员管理

## 🎯 使用场景

### 1. 前端菜单渲染
```javascript
// 获取用户菜单并渲染
const menus = await api.get('/api/v1/user/menus?platform=admin');
renderMenus(menus.data);
```

### 2. 按钮权限控制
```javascript
// 检查是否有创建用户的权限
const canCreateUser = await api.post('/api/v1/user/check-permission', {
  resource: 'user',
  action: 'create'
});
```

### 3. 管理员配置权限
```javascript
// 为角色分配权限
await api.post('/api/v1/admin/permissions/roles/assign-permissions', {
  role_id: 2,
  permission_keys: ['user:list', 'user:create', 'user:update']
});
```

### 4. 管理员配置菜单
```javascript
// 创建新菜单
await api.post('/api/v1/admin/menus', {
  parent_id: 0,
  name: '新模块',
  path: '/new-module',
  permission_key: 'new_module:access',
  platform: 'admin'
});
```

## 🔒 安全特性

1. **分层权限控制**：系统管理 > 功能权限 > 业务权限
2. **平台隔离**：管理端和小程序端权限分离
3. **角色继承**：管理员自动拥有所有权限
4. **实时验证**：基于Casbin的实时权限验证
5. **声明式配置**：路由定义即权限文档

这套完整的权限系统为企业级应用提供了全面的权限控制解决方案！🚀
