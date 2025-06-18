# 用户菜单API接口

## 🎯 概述

用户菜单API提供了用户获取自己可访问菜单、权限和角色的接口，用于前端动态渲染菜单和权限控制。

## 🔧 路由配置

```go
package router

import (
    "github.com/gin-gonic/gin"
    "github.com/limitcool/starter/internal/handler"
    "github.com/limitcool/starter/internal/middleware"
    "github.com/limitcool/starter/internal/pkg/permission"
)

func SetupUserMenuRoutes(
    router *gin.Engine,
    permissionService *permission.Service,
) {
    // 创建用户菜单处理器
    userMenuHandler := handler.NewUserMenuHandler(permissionService)

    // API v1 路由组
    v1 := router.Group("/api/v1")
    
    // 需要认证的路由组
    auth := v1.Group("")
    auth.Use(middleware.JWTMiddleware()) // JWT认证中间件

    // 用户个人信息相关 API
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
}
```

## 📋 API接口列表

| 方法 | 路径 | 功能 | 权限要求 |
|------|------|------|----------|
| GET | `/api/v1/user/menus` | 获取我的菜单 | 仅需登录 |
| GET | `/api/v1/user/permissions` | 获取我的权限 | 仅需登录 |
| GET | `/api/v1/user/roles` | 获取我的角色 | 仅需登录 |
| POST | `/api/v1/user/check-permission` | 检查我的权限 | 仅需登录 |

## 🚀 接口详情

### 1. 获取我的菜单

**请求：**
```bash
GET /api/v1/user/menus?platform=admin
Authorization: Bearer <token>
```

**查询参数：**
- `platform`: 平台类型 (可选，默认: admin)
  - `admin`: 管理端
  - `coach_mp`: 教练小程序端

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
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
      "created_at": "2024-01-01 10:00:00",
      "updated_at": "2024-01-01 10:00:00",
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
      "created_at": "2024-01-01 10:00:00",
      "updated_at": "2024-01-01 10:00:00",
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
          "created_at": "2024-01-01 10:00:00",
          "updated_at": "2024-01-01 10:00:00",
          "children": []
        }
      ]
    }
  ]
}
```

### 2. 获取我的权限

**请求：**
```bash
GET /api/v1/user/permissions
Authorization: Bearer <token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "parent_id": 0,
      "name": "系统管理",
      "key": "sys",
      "type": "MENU",
      "created_at": "2024-01-01 10:00:00",
      "updated_at": "2024-01-01 10:00:00"
    },
    {
      "id": 2,
      "parent_id": 1,
      "name": "查看用户列表",
      "key": "user:list",
      "type": "API",
      "created_at": "2024-01-01 10:00:00",
      "updated_at": "2024-01-01 10:00:00"
    }
  ]
}
```

### 3. 获取我的角色

**请求：**
```bash
GET /api/v1/user/roles
Authorization: Bearer <token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "超级管理员",
      "key": "admin",
      "description": "系统管理员，拥有所有权限",
      "status": 1,
      "created_at": "2024-01-01 10:00:00",
      "updated_at": "2024-01-01 10:00:00"
    },
    {
      "id": 2,
      "name": "教练",
      "key": "coach",
      "description": "管理课程和学员",
      "status": 1,
      "created_at": "2024-01-01 10:00:00",
      "updated_at": "2024-01-01 10:00:00"
    }
  ]
}
```

### 4. 检查我的权限

**请求：**
```bash
POST /api/v1/user/check-permission
Authorization: Bearer <token>
Content-Type: application/json

{
  "resource": "user",
  "action": "create"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "has_permission": true
  }
}
```

## 🔍 权限逻辑

### 1. 菜单过滤逻辑

```
1. 获取用户信息
2. 如果是管理员 → 返回指定平台的所有可见菜单
3. 如果是普通用户：
   a. 获取指定平台的所有可见菜单
   b. 遍历每个菜单：
      - 如果菜单没有权限要求(permission_key为空) → 允许访问
      - 如果菜单有权限要求 → 检查用户是否有该权限
   c. 构建菜单树结构
4. 返回过滤后的菜单树
```

### 2. 权限检查逻辑

```
1. 从JWT token中获取用户ID
2. 调用Casbin服务检查用户权限
3. 返回权限检查结果
```

## 🎯 前端使用示例

### 1. Vue.js 菜单渲染

```javascript
// 获取用户菜单
async function getUserMenus(platform = 'admin') {
  try {
    const response = await api.get(`/api/v1/user/menus?platform=${platform}`);
    return response.data.data;
  } catch (error) {
    console.error('获取用户菜单失败:', error);
    return [];
  }
}

// 渲染菜单组件
export default {
  data() {
    return {
      menus: []
    };
  },
  async mounted() {
    this.menus = await getUserMenus('admin');
  },
  methods: {
    renderMenu(menus) {
      return menus.map(menu => ({
        id: menu.id,
        title: menu.name,
        icon: menu.icon,
        path: menu.path,
        component: menu.component,
        children: menu.children ? this.renderMenu(menu.children) : []
      }));
    }
  }
};
```

### 2. React 权限控制

```javascript
// 权限检查Hook
import { useState, useEffect } from 'react';

function usePermission(resource, action) {
  const [hasPermission, setHasPermission] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function checkPermission() {
      try {
        const response = await api.post('/api/v1/user/check-permission', {
          resource,
          action
        });
        setHasPermission(response.data.data.has_permission);
      } catch (error) {
        console.error('权限检查失败:', error);
        setHasPermission(false);
      } finally {
        setLoading(false);
      }
    }

    checkPermission();
  }, [resource, action]);

  return { hasPermission, loading };
}

// 权限控制组件
function PermissionButton({ resource, action, children, ...props }) {
  const { hasPermission, loading } = usePermission(resource, action);

  if (loading) return null;
  if (!hasPermission) return null;

  return <button {...props}>{children}</button>;
}

// 使用示例
function UserManagement() {
  return (
    <div>
      <h1>用户管理</h1>
      <PermissionButton resource="user" action="create">
        创建用户
      </PermissionButton>
      <PermissionButton resource="user" action="delete">
        删除用户
      </PermissionButton>
    </div>
  );
}
```

## 🔒 安全特性

1. **JWT认证**：所有接口都需要有效的JWT token
2. **用户隔离**：用户只能获取自己的菜单和权限
3. **实时权限**：权限检查基于Casbin实时策略
4. **平台隔离**：不同平台的菜单相互隔离
5. **管理员特权**：管理员自动拥有所有权限

## 📁 核心文件

- `internal/handler/user_menu_handler.go` - 用户菜单处理器
- `internal/pkg/permission/service.go` - 权限服务（GetUserMenusByPlatform方法）
- `docs/USER_MENU_API.md` - API文档

这套用户菜单API为前端提供了完整的权限控制和菜单渲染支持！🚀
