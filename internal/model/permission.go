package model

// Permission 权限实体
type Permission struct {
	BaseModel

	Name        string `json:"name" gorm:"size:50;not null;comment:权限名称"`
	Code        string `json:"code" gorm:"size:100;not null;unique;comment:权限编码"`
	Description string `json:"description" gorm:"size:200;comment:权限描述"`
	Type        int8   `json:"type" gorm:"default:1;comment:权限类型(1:菜单,2:操作,3:API)"`
	Status      int8   `json:"status" gorm:"default:1;comment:状态(0:禁用,1:正常)"`

	// 关联的菜单
	MenuID uint  `json:"menu_id" gorm:"comment:关联菜单ID"`
	Menu   *Menu `json:"menu" gorm:"foreignKey:MenuID"`

	// API相关
	ApiPath   string `json:"api_path" gorm:"size:200;comment:API路径"`
	ApiMethod string `json:"api_method" gorm:"size:10;comment:API方法(GET,POST,PUT,DELETE)"`
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
