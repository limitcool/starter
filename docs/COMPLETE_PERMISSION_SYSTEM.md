# 完整的Casbin权限系统实现

## 🎯 系统概述

我们实现了一个基于Casbin的完整权限系统，遵循最佳实践，实现了真正的解耦设计。

## 📊 核心架构

### 设计原则：角色 → 权限 → 菜单 (通过权限标识连接)

```
用户 ←→ 角色 ←→ 权限标识 ←→ 菜单
     (Casbin)    (权限字典)    (permission_key)
```

### 数据流向

1. **权限配置**：管理员在界面选择权限 → permissions表展示 → 转换为permission_key → Casbin API存储策略
2. **菜单渲染**：用户登录 → Casbin获取权限 → 过滤菜单permission_key → 返回可访问菜单  
3. **API保护**：请求API → 权限中间件 → Casbin验证permission_key → 允许/拒绝

## 🗄️ 数据模型

### 1. 核心表结构

```sql
-- 用户表
CREATE TABLE users (
    id bigint PRIMARY KEY AUTO_INCREMENT,
    uuid varchar(36) UNIQUE,
    username varchar(50) UNIQUE,
    password varchar(255),
    nickname varchar(50),
    phone varchar(20) UNIQUE,
    openid varchar(128) UNIQUE,
    status tinyint(1) DEFAULT 1,
    is_admin boolean DEFAULT false,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

-- 角色表
CREATE TABLE roles (
    id bigint PRIMARY KEY AUTO_INCREMENT,
    name varchar(50) UNIQUE,        -- 角色名 (如: 超级管理员)
    `key` varchar(50) UNIQUE,       -- 角色标识 (如: admin, coach)
    description varchar(255),       -- 角色描述
    status tinyint(1) DEFAULT 1,    -- 角色状态
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

-- 权限字典表
CREATE TABLE permissions (
    id bigint PRIMARY KEY AUTO_INCREMENT,
    parent_id bigint DEFAULT 0,     -- 父权限ID (用于分组)
    name varchar(50),               -- 权限名 (如: 查看会员列表)
    `key` varchar(100) UNIQUE,      -- 权限标识 (如: member:list)
    type varchar(20),               -- 权限类型 (MENU:菜单权限, BUTTON:按钮权限, API:接口权限)
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

-- 菜单表
CREATE TABLE menus (
    id bigint PRIMARY KEY AUTO_INCREMENT,
    parent_id bigint DEFAULT 0,     -- 父菜单ID
    name varchar(50),               -- 菜单名 (如: 用户管理)
    path varchar(255),              -- 前端路由路径 (如: /user)
    component varchar(255),         -- 前端组件路径
    icon varchar(50),               -- 菜单图标
    sort_order int DEFAULT 0,       -- 显示排序
    is_visible boolean DEFAULT true, -- 是否可见
    permission_key varchar(100),    -- 核心字段：访问此菜单所需的权限标识
    platform varchar(20) DEFAULT 'admin', -- 所属平台 (admin:管理端, coach_mp:教练小程序端)
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

-- Casbin策略表 (自动创建)
CREATE TABLE casbin_rule (
    id bigint PRIMARY KEY AUTO_INCREMENT,
    ptype varchar(100),             -- 策略类型 (p:权限策略, g:角色关联)
    v0 varchar(100),                -- Subject (用户/角色)
    v1 varchar(100),                -- Object (资源/权限标识)
    v2 varchar(100),                -- Action (操作)
    v3 varchar(100),
    v4 varchar(100),
    v5 varchar(100),
    UNIQUE KEY idx_casbin_rule (ptype,v0,v1,v2,v3,v4,v5)
);
```

### 2. 示例数据

**角色数据：**
```sql
INSERT INTO roles (name, `key`, description) VALUES
('超级管理员', 'admin', '系统管理员，拥有所有权限'),
('教练', 'coach', '管理课程和学员'),
('销售', 'sales', '管理会员和线索'),
('会员', 'member', '预约课程');
```

**权限字典数据：**
```sql
INSERT INTO permissions (parent_id, name, `key`, type) VALUES
(0, '系统管理', 'sys', 'MENU'),
(1, '用户管理', 'sys:user', 'MENU'),
(2, '查看用户列表', 'user:list', 'API'),
(2, '创建用户', 'user:create', 'API'),
(0, '会员管理', 'member_manage', 'MENU'),
(5, '查看会员列表', 'member:list', 'API'),
(5, '编辑会员信息', 'member:edit', 'API'),
(0, '课程管理', 'course_manage', 'MENU'),
(8, '查看课程列表', 'course:list', 'API'),
(8, '创建课程', 'course:create', 'API'),
(0, '我的学员(小程序)', 'mp_student:list', 'API');
```

**菜单数据：**
```sql
INSERT INTO menus (parent_id, name, path, component, icon, sort_order, permission_key, platform) VALUES
(0, '仪表盘', '/dashboard', 'Dashboard', 'dashboard', 1, '', 'admin'),
(0, '系统管理', '/system', 'Layout', 'system', 2, 'sys', 'admin'),
(2, '用户管理', '/system/user', 'system/User', 'user', 1, 'sys:user', 'admin'),
(0, '会员管理', '/member', 'member/Index', 'member', 3, 'member_manage', 'admin'),
(0, '课程管理', '/course', 'course/Index', 'course', 4, 'course_manage', 'admin'),
(0, '我的学员', '/mp/students', 'mp/Students', 'student', 1, 'mp_student:list', 'coach_mp');
```

**Casbin策略数据：**
```sql
-- 权限策略 (p策略)
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES 
('p', 'admin', '*', '*'),              -- admin角色拥有所有权限
('p', 'coach', 'course', 'list'),      -- coach角色可以查看课程
('p', 'coach', 'course', 'create'),    -- coach角色可以创建课程
('p', 'coach', 'member', 'list'),      -- coach角色可以查看会员
('p', 'coach', 'mp_student', 'list'),  -- coach角色可以查看小程序学员
('p', 'sales', 'member', 'list'),      -- sales角色可以查看会员
('p', 'sales', 'member', 'edit');      -- sales角色可以编辑会员

-- 用户角色关联 (g策略)
INSERT INTO casbin_rule (ptype, v0, v1) VALUES
('g', 'user:1', 'admin'),              -- 用户1是admin角色
('g', 'user:2', 'coach'),              -- 用户2是coach角色
('g', 'user:3', 'sales');              -- 用户3是sales角色
```

## 🔧 核心组件

### 1. Casbin服务 (`internal/pkg/casbin/casbin.go`)

```go
// 权限验证
func (s *Service) Enforce(ctx context.Context, user, resource, action string) (bool, error)

// 添加权限策略
func (s *Service) AddPolicy(ctx context.Context, role, resource, action string) error

// 用户角色管理
func (s *Service) AddRoleForUser(ctx context.Context, user, role string) error
func (s *Service) GetRolesForUser(ctx context.Context, user string) ([]string, error)
```

### 2. 权限服务 (`internal/pkg/permission/service.go`)

```go
// 权限检查
func (s *Service) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error)

// 角色权限分配
func (s *Service) AssignPermissionsToRole(ctx context.Context, roleID uint, permissionKeys []string) error

// 用户角色分配
func (s *Service) AssignRolesToUser(ctx context.Context, userID int64, roleKeys []string) error

// 获取用户菜单
func (s *Service) GetUserMenus(ctx context.Context, userID int64) ([]model.Menu, error)
```

### 3. 权限中间件 (`internal/middleware/permission.go`)

```go
// 基于权限Key的中间件 (推荐)
func (m *PermissionMiddleware) RequirePermissionKey(permissionKey string) gin.HandlerFunc

// 基于资源和操作的中间件
func (m *PermissionMiddleware) RequirePermission(resource, action string) gin.HandlerFunc

// 管理员权限中间件
func (m *PermissionMiddleware) RequireAdmin() gin.HandlerFunc

// 角色权限中间件
func (m *PermissionMiddleware) RequireRole(roleName string) gin.HandlerFunc
```

## 🚀 使用示例

### 1. 路由权限控制

```go
// 创建权限中间件
permissionMiddleware := middleware.NewPermissionMiddleware(permissionService)

// 用户管理 API
v1.GET("/users", permissionMiddleware.RequirePermissionKey("user:list"), userHandler.List)
v1.POST("/users", permissionMiddleware.RequirePermissionKey("user:create"), userHandler.Create)
v1.PUT("/users/:id", permissionMiddleware.RequirePermissionKey("user:update"), userHandler.Update)
v1.DELETE("/users/:id", permissionMiddleware.RequirePermissionKey("user:delete"), userHandler.Delete)

// 管理员专用接口
v1.GET("/admin/config", permissionMiddleware.RequireAdmin(), adminHandler.GetConfig)

// 角色专用接口
v1.GET("/coach/students", permissionMiddleware.RequireRole("教练"), coachHandler.GetStudents)
```

### 2. 权限分配

```go
// 为角色分配权限
err := permissionService.AssignPermissionsToRole(ctx, roleID, []string{
    "user:create", "user:list", "user:update"
})

// 为用户分配角色
err := permissionService.AssignRolesToUser(ctx, userID, []string{
    "coach", "sales"
})
```

### 3. 菜单渲染

```go
// 获取用户可访问的菜单
menus, err := permissionService.GetUserMenus(ctx, userID)
// 返回根据用户权限过滤后的菜单树
```

## 🎉 系统优势

### 1. 真正解耦
- **消除冗余关联表**：不需要user_roles、role_permissions、role_menus、menu_permissions等表
- **菜单权限解耦**：菜单通过permission_key关联权限，不直接绑定角色
- **权限验证统一**：所有权限验证都通过Casbin进行

### 2. 职责分离
- **permissions表**：给"人"看的权限字典，用于管理界面展示
- **casbin_rule表**：给"机器"用的权限策略，用于运行时验证
- **menus表**：通过permission_key桥接权限，实现动态菜单

### 3. 高度灵活
- **声明式权限控制**：路由定义即权限文档
- **多种权限控制方式**：支持权限Key、角色、管理员等多种控制方式
- **平台权限区分**：支持管理端和小程序端的权限区分

### 4. 易于维护
- **权限变更简单**：只需修改中间件参数
- **业务逻辑解耦**：业务代码无需关心权限检查
- **统一错误处理**：权限错误统一在中间件中处理

## 📁 核心文件

- `internal/model/role.go` - 数据模型定义
- `internal/pkg/casbin/casbin.go` - Casbin服务封装
- `internal/pkg/permission/service.go` - 权限业务逻辑
- `internal/middleware/permission.go` - 权限中间件
- `internal/migration/20250618_create_permission_tables.go` - 数据库迁移
- `docs/PERMISSION_MIDDLEWARE_USAGE.md` - 中间件使用指南

这个权限系统真正实现了Casbin的最佳实践，提供了完整、灵活、易维护的权限控制解决方案！🎯
