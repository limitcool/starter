package model

import (
	"github.com/limitcool/starter/internal/pkg/enum"
)

// Permission 权限实体
type Permission struct {
	BaseModel

	Name        string              `json:"name" gorm:"size:50;not null;comment:权限名称"`
	Code        string              `json:"code" gorm:"size:100;not null;unique;index;comment:权限编码"`
	Type        enum.PermissionType `json:"type" gorm:"default:0;index;comment:权限类型(0:菜单,1:按钮,2:API)"`
	Description string              `json:"description" gorm:"size:200;comment:权限描述"`
	Enabled     bool                `json:"enabled" gorm:"default:true;comment:是否启用"`
	Remark      string              `json:"remark" gorm:"size:500;comment:备注"`

	// 关联字段
	MenuID   uint `json:"menu_id" gorm:"default:0;index;comment:关联菜单ID"`
	ButtonID uint `json:"button_id" gorm:"default:0;index;comment:关联按钮ID"`
	APIID    uint `json:"api_id" gorm:"default:0;index;comment:关联接口ID"`

	// 关联对象
	Menu   *Menu       `json:"menu" gorm:"foreignKey:MenuID"`               // 所属菜单
	Button *MenuButton `json:"button" gorm:"foreignKey:ButtonID"`           // 所属按钮
	API    *API        `json:"api" gorm:"foreignKey:APIID"`                 // 所属接口
	Roles  []*Role     `json:"roles" gorm:"many2many:sys_role_permission;"` // 关联的角色

	// 非数据库字段
	RoleIDs []uint `json:"role_ids" gorm:"-"` // 角色ID列表，不映射到数据库
}

// 表名
func (Permission) TableName() string {
	return "sys_permission"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	BaseModel

	RoleID       uint `json:"role_id" gorm:"not null;index;comment:角色ID"`
	PermissionID uint `json:"permission_id" gorm:"not null;index;comment:权限ID"`
}

// 表名
func (RolePermission) TableName() string {
	return "sys_role_permission"
}

// 以下方法已移动到 repository/permission_repo.go
// Create
// Update
// Delete
// GetByID
// GetAll
// GetByRoleID
// GetByUserID
// BatchCreate
// DeleteByRoleID
