# 权限和菜单CRUD路由配置

## 🎯 概述

我们已经实现了完整的权限和菜单CRUD操作，包括增删改查、树形结构展示等功能。

## 🔧 路由配置示例

```go
package router

import (
    "github.com/gin-gonic/gin"
    "github.com/limitcool/starter/internal/handler"
    "github.com/limitcool/starter/internal/middleware"
    "github.com/limitcool/starter/internal/model"
    "github.com/limitcool/starter/internal/pkg/permission"
)

func SetupCRUDRoutes(
    router *gin.Engine,
    permissionService *permission.Service,
    permissionRepo *model.PermissionRepo,
    menuRepo *model.MenuRepo,
) {
    // 创建权限中间件
    permissionMiddleware := middleware.NewPermissionMiddleware(permissionService)
    
    // 创建CRUD处理器
    permissionCRUDHandler := handler.NewPermissionCRUDHandler(permissionRepo)
    menuCRUDHandler := handler.NewMenuCRUDHandler(menuRepo)

    // API v1 路由组
    v1 := router.Group("/api/v1")
    
    // 需要认证的路由组
    auth := v1.Group("")
    auth.Use(middleware.JWTMiddleware()) // JWT认证中间件

    // 权限管理 API
    permissionGroup := auth.Group("/admin/permissions")
    {
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

    // 菜单管理 API
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
}
```

## 📋 API接口列表

### 权限管理接口

| 方法 | 路径 | 功能 | 权限要求 |
|------|------|------|----------|
| GET | `/api/v1/admin/permissions` | 获取权限列表 | `permission:list` |
| GET | `/api/v1/admin/permissions/tree` | 获取权限树 | `permission:list` |
| GET | `/api/v1/admin/permissions/:id` | 获取权限详情 | `permission:view` |
| POST | `/api/v1/admin/permissions` | 创建权限 | `permission:create` |
| PUT | `/api/v1/admin/permissions/:id` | 更新权限 | `permission:update` |
| DELETE | `/api/v1/admin/permissions/:id` | 删除权限 | `permission:delete` |

### 菜单管理接口

| 方法 | 路径 | 功能 | 权限要求 |
|------|------|------|----------|
| GET | `/api/v1/admin/menus` | 获取菜单列表 | `menu:list` |
| GET | `/api/v1/admin/menus/tree` | 获取菜单树 | `menu:list` |
| GET | `/api/v1/admin/menus/:id` | 获取菜单详情 | `menu:view` |
| POST | `/api/v1/admin/menus` | 创建菜单 | `menu:create` |
| PUT | `/api/v1/admin/menus/:id` | 更新菜单 | `menu:update` |
| DELETE | `/api/v1/admin/menus/:id` | 删除菜单 | `menu:delete` |
| PUT | `/api/v1/admin/menus/sort` | 更新菜单排序 | `menu:update` |

## 🚀 请求示例

### 1. 创建权限

```bash
POST /api/v1/admin/permissions
Content-Type: application/json
Authorization: Bearer <token>

{
  "parent_id": 0,
  "name": "查看用户列表",
  "key": "user:list",
  "type": "API"
}
```

### 2. 获取权限树

```bash
GET /api/v1/admin/permissions/tree
Authorization: Bearer <token>

Response:
[
  {
    "id": 1,
    "parent_id": 0,
    "name": "系统管理",
    "key": "sys",
    "type": "MENU",
    "children": [
      {
        "id": 2,
        "parent_id": 1,
        "name": "用户管理",
        "key": "sys:user",
        "type": "MENU",
        "children": [
          {
            "id": 3,
            "parent_id": 2,
            "name": "查看用户列表",
            "key": "user:list",
            "type": "API",
            "children": []
          }
        ]
      }
    ]
  }
]
```

### 3. 创建菜单

```bash
POST /api/v1/admin/menus
Content-Type: application/json
Authorization: Bearer <token>

{
  "parent_id": 0,
  "name": "用户管理",
  "path": "/system/user",
  "component": "system/User",
  "icon": "user",
  "sort_order": 1,
  "is_visible": true,
  "permission_key": "sys:user",
  "platform": "admin"
}
```

### 4. 获取菜单树

```bash
GET /api/v1/admin/menus/tree?platform=admin
Authorization: Bearer <token>

Response:
[
  {
    "id": 1,
    "parent_id": 0,
    "name": "仪表盘",
    "path": "/dashboard",
    "component": "Dashboard",
    "icon": "dashboard",
    "sort_order": 1,
    "is_visible": true,
    "permission_key": "",
    "platform": "admin",
    "children": []
  },
  {
    "id": 2,
    "parent_id": 0,
    "name": "系统管理",
    "path": "/system",
    "component": "Layout",
    "icon": "system",
    "sort_order": 2,
    "is_visible": true,
    "permission_key": "sys",
    "platform": "admin",
    "children": [
      {
        "id": 3,
        "parent_id": 2,
        "name": "用户管理",
        "path": "/system/user",
        "component": "system/User",
        "icon": "user",
        "sort_order": 1,
        "is_visible": true,
        "permission_key": "sys:user",
        "platform": "admin",
        "children": []
      }
    ]
  }
]
```

### 5. 更新菜单排序

```bash
PUT /api/v1/admin/menus/sort
Content-Type: application/json
Authorization: Bearer <token>

{
  "menu_sorts": [
    {"id": 1, "sort_order": 1},
    {"id": 2, "sort_order": 2},
    {"id": 3, "sort_order": 3}
  ]
}
```

## 🔍 查询参数

### 权限列表查询参数

- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 10)
- `parent_id`: 父权限ID (可选)

### 菜单列表查询参数

- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 10)
- `parent_id`: 父菜单ID (可选)
- `platform`: 平台 (默认: admin, 可选值: admin, coach_mp)

## 🎯 特性

### 1. 树形结构支持
- 权限和菜单都支持无限级树形结构
- 提供专门的树形接口，返回完整的层级关系

### 2. 分页支持
- 列表接口支持分页查询
- 返回总数和分页信息

### 3. 权限控制
- 所有接口都有相应的权限控制
- 使用声明式权限中间件

### 4. 数据验证
- 完整的请求参数验证
- 业务逻辑验证（如删除前检查子项）

### 5. 平台区分
- 菜单支持平台区分（管理端、小程序端）
- 可以为不同平台配置不同的菜单

## 📁 核心文件

- `internal/handler/permission_crud_handler.go` - 权限CRUD处理器
- `internal/handler/menu_crud_handler.go` - 菜单CRUD处理器
- `internal/dto/permission.go` - 请求响应DTO定义
- `internal/model/role_repo.go` - Repository层实现

这套CRUD系统为权限和菜单管理提供了完整的后台管理功能！🚀
