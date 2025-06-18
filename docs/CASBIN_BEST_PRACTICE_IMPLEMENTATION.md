# 基于Casbin的权限系统最佳实践实现

## 🎯 设计理念

感谢您的深入质疑！我们现在实现了一个真正基于Casbin的权限系统，完全遵循最佳实践：

### 核心架构：权限字典 + Casbin策略 + 解耦菜单

```
用户 ←→ 角色 ←→ 权限标识 ←→ 菜单
     (Casbin)    (权限字典)    (permission_key)
```

## 📊 数据模型设计

### 1. User表 - 用户模型
```go
type User struct {
    gorm.Model
    UUID     string `gorm:"type:varchar(36);uniqueIndex"`
    Username string `gorm:"type:varchar(50);uniqueIndex"`
    Password string `gorm:"type:varchar(255)"`
    Nickname string `gorm:"type:varchar(50)"`
    Phone    string `gorm:"type:varchar(20);uniqueIndex"`
    OpenID   string `gorm:"type:varchar(128);uniqueIndex"`
    Status   uint8  `gorm:"type:tinyint(1);default:1"`
    IsAdmin  bool   `gorm:"default:false"`
    
    // 注意：用户角色关系由Casbin管理，不再需要Roles字段
}
```

### 2. Role表 - 角色模型
```go
type Role struct {
    gorm.Model
    Name        string `gorm:"type:varchar(50);uniqueIndex"`        // 角色名 (如: 超级管理员)
    Key         string `gorm:"type:varchar(50);uniqueIndex"`        // 角色唯一标识 (如: admin, coach)
    Description string `gorm:"type:varchar(255)"`                  // 角色描述
    Status      uint8  `gorm:"type:tinyint(1);default:1"`          // 角色状态
    
    // 注意：角色权限关系由Casbin管理，不再需要Permissions字段
}
```

### 3. Permission表 - 权限字典
```go
type Permission struct {
    gorm.Model
    ParentID uint   `gorm:"default:0;index"`                       // 父权限ID (用于分组)
    Name     string `gorm:"type:varchar(50)"`                      // 权限名 (如: 查看会员列表)
    Key      string `gorm:"type:varchar(100);uniqueIndex"`         // 权限唯一标识 (如: member:list)
    Type     string `gorm:"type:varchar(20)"`                      // 权限类型 (MENU:菜单权限, BUTTON:按钮权限, API:接口权限)
}
```

### 4. Menu表 - 菜单模型
```go
type Menu struct {
    gorm.Model
    ParentID      uint   `gorm:"default:0;index"`                  // 父菜单ID
    Name          string `gorm:"type:varchar(50)"`                 // 菜单名 (如: 用户管理)
    Path          string `gorm:"type:varchar(255)"`                // 前端路由路径 (如: /user)
    Component     string `gorm:"type:varchar(255)"`                // 前端组件路径
    Icon          string `gorm:"type:varchar(50)"`                 // 菜单图标
    SortOrder     int    `gorm:"type:int;default:0"`               // 显示排序
    IsVisible     bool   `gorm:"type:tinyint(1);default:1"`        // 是否可见
    PermissionKey string `gorm:"type:varchar(100)"`                // 核心字段：访问此菜单所需的权限标识
    Platform      string `gorm:"type:varchar(20);default:'admin'"` // 所属平台 (admin:管理端, coach_mp:教练小程序端)
    Children      []Menu `gorm:"foreignKey:ParentID"`              // 子菜单
}
```

### 5. Casbin策略表 - 权限验证
```sql
-- casbin_rule表存储实际的权限策略
CREATE TABLE `casbin_rule` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) DEFAULT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 🔧 核心工作流程

### 1. 权限分配流程
```
管理员在界面选择权限 → 读取permissions表展示 → 选择后转换为permission_key → 调用Casbin API存储策略
```

**示例数据：**
```sql
-- 权限策略 (p策略)
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES 
('p', 'admin', '*', '*'),              -- admin角色拥有所有权限
('p', 'coach', 'course', 'list'),      -- coach角色可以查看课程
('p', 'coach', 'course', 'create'),    -- coach角色可以创建课程
('p', 'coach', 'member', 'list'),      -- coach角色可以查看会员
('p', 'sales', 'member', 'list'),      -- sales角色可以查看会员
('p', 'sales', 'member', 'edit');      -- sales角色可以编辑会员

-- 用户角色关联 (g策略)
INSERT INTO casbin_rule (ptype, v0, v1) VALUES
('g', 'user:1', 'admin'),              -- 用户1是admin角色
('g', 'user:2', 'coach'),              -- 用户2是coach角色
('g', 'user:3', 'sales');              -- 用户3是sales角色
```

### 2. 菜单渲染流程
```
用户登录 → Casbin获取用户权限列表 → 遍历菜单检查permission_key → 返回可访问菜单
```

**核心代码：**
```go
func (s *Service) GetUserMenus(ctx context.Context, userID int64) ([]model.Menu, error) {
    allMenus, _ := s.menuRepo.GetEnabledMenus(ctx)
    var accessibleMenus []model.Menu
    
    for _, menu := range allMenus {
        if menu.PermissionKey == "" {
            // 无权限要求的菜单
            accessibleMenus = append(accessibleMenus, menu)
            continue
        }
        
        // 检查用户权限
        parts := strings.Split(menu.PermissionKey, ":")
        if len(parts) == 2 {
            hasPermission, _ := s.CheckPermission(ctx, userID, parts[0], parts[1])
            if hasPermission {
                accessibleMenus = append(accessibleMenus, menu)
            }
        }
    }
    
    return s.menuRepo.BuildMenuTree(accessibleMenus, 0), nil
}
```

### 3. API权限验证流程
```
请求API → 中间件调用Casbin.Enforce(user, resource, action) → 允许/拒绝访问
```

## 🎉 设计优势

### 1. 职责分离
- **permissions表**：给"人"看的权限字典，用于管理界面展示
- **casbin_rule表**：给"机器"用的权限策略，用于运行时验证
- **menus表**：通过permission_key桥接权限，实现动态菜单

### 2. 真正解耦
- **消除冗余关联表**：不再需要user_roles、role_permissions、role_menus、menu_permissions等表
- **菜单权限解耦**：菜单通过permission_key关联权限，不直接绑定角色
- **权限验证统一**：所有权限验证都通过Casbin进行

### 3. 灵活性
- **新增权限**：只需在permissions表添加记录，然后通过管理界面分配
- **菜单权限控制**：修改menu.permission_key即可改变菜单访问权限
- **角色权限分配**：通过Casbin API动态管理，支持复杂的权限继承

### 4. 可维护性
- **权限、角色、菜单各司其职**
- **修改一处不影响全局**
- **逻辑清晰，易于理解和扩展**

## 🚀 API接口示例

### 分配角色权限
```json
POST /api/admin/permissions/roles/assign-permissions
{
  "role_id": 1,
  "permission_keys": ["user:create", "user:list", "user:update"]
}
```

### 分配用户角色
```json
POST /api/admin/permissions/assign-user-roles
{
  "user_id": 1,
  "role_keys": ["admin", "coach"]
}
```

### 获取用户菜单
```json
GET /api/user/menus
Response: [
  {
    "id": 1,
    "name": "用户管理",
    "path": "/system/user",
    "permission_key": "sys:user",
    "platform": "admin",
    "children": []
  }
]
```

## 📝 总结

这个重新设计的权限系统完全遵循了您提到的最佳实践：

1. **Casbin负责权限验证**：所有权限检查都通过Casbin进行
2. **permissions表作为权限字典**：用于管理界面展示和权限分配
3. **菜单通过permission_key关联权限**：实现了真正的解耦设计
4. **消除了所有冗余关联表**：简化了数据模型
5. **权限分配通过业务逻辑连接**：从permissions表获取Key，调用Casbin API存储策略

感谢您的深入思考和质疑，这让我们实现了一个真正优秀的权限系统！🎯
