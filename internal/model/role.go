package model

// 角色实体
type Role struct {
	BaseModel

	Name        string `json:"name" gorm:"size:50;not null;unique;comment:角色名称"`
	Code        string `json:"code" gorm:"size:50;not null;unique;comment:角色编码"`
	Status      int8   `json:"status" gorm:"default:1;comment:状态(0:禁用,1:正常)"`
	Sort        int    `json:"sort" gorm:"default:0;comment:排序"`
	Description string `json:"description" gorm:"size:100;comment:角色描述"`

	// 关联的菜单ID列表(用于前端传递)
	MenuIDs []uint `json:"menu_ids" gorm:"-"`
}

// 表名
func (Role) TableName() string {
	return "sys_role"
}

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

// 用户角色关联表
type UserRole struct {
	BaseModel

	UserID uint `json:"user_id" gorm:"not null;comment:用户ID"`
	RoleID uint `json:"role_id" gorm:"not null;comment:角色ID"`
}

// 表名
func (UserRole) TableName() string {
	return "sys_user_role"
}
