package model

import (
	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(50);uniqueIndex;comment:角色名 (如: 超级管理员)" json:"name"`
	Key         string `gorm:"type:varchar(50);uniqueIndex;comment:角色唯一标识 (如: admin, coach)" json:"key"` // Casbin 中的 role/subject
	Description string `gorm:"type:varchar(255);comment:角色描述" json:"description"`
	Status      uint8  `gorm:"type:tinyint(1);default:1;comment:角色状态 (1:正常 2:禁用)" json:"status"`

	// 注意：这里不再需要 Permissions、Users、Menus 字段，因为关系由 Casbin 管理
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限字典表 (用于UI展示和管理)
type Permission struct {
	gorm.Model
	ParentID uint   `gorm:"default:0;index;comment:父权限ID (用于分组)" json:"parent_id"`
	Name     string `gorm:"type:varchar(50);comment:权限名 (如: 查看会员列表)" json:"name"`
	// Key 是核心，它将和 Casbin 的 obj 对应
	Key  string `gorm:"type:varchar(100);uniqueIndex;comment:权限唯一标识 (如: member:list)" json:"key"`
	Type string `gorm:"type:varchar(20);comment:权限类型 (MENU:菜单权限, BUTTON:按钮权限, API:接口权限)" json:"type"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// Menu 菜单模型
type Menu struct {
	gorm.Model
	ParentID  uint   `gorm:"default:0;index;comment:父菜单ID" json:"parent_id"`
	Name      string `gorm:"type:varchar(50);comment:菜单名 (如: 用户管理)" json:"name"`
	Path      string `gorm:"type:varchar(255);comment:前端路由路径 (如: /user)" json:"path"`
	Component string `gorm:"type:varchar(255);comment:前端组件路径 (如: /views/user/index.vue)" json:"component"`
	Icon      string `gorm:"type:varchar(50);comment:菜单图标" json:"icon"`
	SortOrder int    `gorm:"type:int;default:0;comment:显示排序" json:"sort_order"`
	IsVisible bool   `gorm:"type:tinyint(1);default:1;comment:是否可见 (用于固定菜单，不走权限)" json:"is_visible"`

	// 核心关联字段: 指向 permissions 表的 Key
	PermissionKey string `gorm:"type:varchar(100);comment:访问此菜单所需的权限标识" json:"permission_key"`

	// 平台区分字段
	Platform string `gorm:"type:varchar(20);default:'admin';index;comment:所属平台 (admin:管理端, coach_mp:教练小程序端)" json:"platform"`

	// 关联关系
	Children []Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// 注意：在基于Casbin的设计中，我们不需要以下关联表：
// - UserRole: 用户角色关联由Casbin的g策略管理
// - RolePermission: 角色权限关联由Casbin的p策略管理
// - RoleMenu: 菜单通过permission_key与权限关联，不直接关联角色
// - MenuPermission: 菜单通过permission_key字段关联权限
