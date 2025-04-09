package model

import "github.com/limitcool/starter/internal/pkg/enum"

// Permission 权限实体
type Permission struct {
	BaseModel

	Name     string              `json:"name" gorm:"size:50;not null;comment:权限名称"`
	Code     string              `json:"code" gorm:"size:100;not null;unique;comment:权限编码"`
	Type     enum.PermissionType `json:"type" gorm:"default:0;comment:权限类型(0:菜单,1:操作,2:API)"`
	MenuID   uint                `json:"menu_id" gorm:"default:0;comment:所属菜单ID"`
	ParentID uint                `json:"parent_id" gorm:"default:0;comment:父权限ID"`
	Path     string              `json:"path" gorm:"size:100;comment:权限路径"`
	Method   string              `json:"method" gorm:"size:10;comment:请求方法"`
	Enabled  bool                `json:"enabled" gorm:"default:true;comment:是否启用"`
	Remark   string              `json:"remark" gorm:"size:500;comment:备注"`

	// 关联
	Menu    *Menu   `json:"menu" gorm:"foreignKey:MenuID"`               // 所属菜单
	Roles   []*Role `json:"roles" gorm:"many2many:sys_role_permission;"` // 关联的角色
	RoleIDs []uint  `json:"role_ids" gorm:"-"`                           // 角色ID列表，不映射到数据库
}

// 表名
func (Permission) TableName() string {
	return "sys_permission"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	BaseModel

	RoleID       uint `json:"role_id" gorm:"not null;comment:角色ID"`
	PermissionID uint `json:"permission_id" gorm:"not null;comment:权限ID"`
}

// 表名
func (RolePermission) TableName() string {
	return "sys_role_permission"
}
