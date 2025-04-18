package model

// 角色实体
type Role struct {
	BaseModel

	Name        string `json:"name" gorm:"size:50;not null;unique;comment:角色名称"`
	Code        string `json:"code" gorm:"size:100;not null;unique;comment:角色编码"`
	Enabled     bool   `json:"enabled" gorm:"default:true;comment:是否启用"`
	Sort        int    `json:"sort" gorm:"default:0;comment:排序"`
	Description string `json:"description" gorm:"size:200;comment:角色描述"`
	Remark      string `json:"remark" gorm:"size:500;comment:备注"`

	// 关联
	Users         []*User       `json:"users" gorm:"many2many:sys_user_role;"`             // 关联的用户
	Permissions   []*Permission `json:"permissions" gorm:"many2many:sys_role_permission;"` // 关联的权限
	Menus         []*Menu       `json:"menus" gorm:"many2many:sys_role_menu;"`             // 关联的菜单
	UserIDs       []uint        `json:"user_ids" gorm:"-"`                                 // 用户ID列表，不映射到数据库
	PermissionIDs []uint        `json:"permission_ids" gorm:"-"`                           // 权限ID列表，不映射到数据库
	MenuIDs       []uint        `json:"menu_ids" gorm:"-"`                                 // 菜单ID列表，不映射到数据库
}

// 表名
func (Role) TableName() string {
	return "sys_role"
}

// 以下方法已移动到 repository/role_repo.go
// Create
// Update
// Delete
// GetByID
// GetAll
// IsAssignedToUser
// DeleteRoleMenus

// 角色菜单关联表
type RoleMenu struct {
	BaseModel

	RoleID uint `json:"role_id" gorm:"not null;comment:角色ID"`
	MenuID uint `json:"menu_id" gorm:"not null;comment:菜单ID"`
}

// 表名
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

// 以下方法已移动到 repository/role_repo.go
// BatchCreate
// GetMenuIDsByRoleID

// 用户角色关联表
type UserRole struct {
	BaseModel

	UserID int64 `json:"user_id" gorm:"type:bigint;not null;comment:用户ID"`
	RoleID uint  `json:"role_id" gorm:"not null;comment:角色ID"`
}

// 表名
func (UserRole) TableName() string {
	return "sys_user_role"
}

// 以下方法已移动到 repository/role_repo.go
// BatchCreate
// DeleteByUserID
// GetRoleIDsByUserID
