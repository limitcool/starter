package model

// AdminUserRole 管理员用户角色关联表
type AdminUserRole struct {
	BaseModel

	AdminUserID int64 `json:"admin_user_id" gorm:"type:bigint;not null;index;comment:管理员用户ID"`
	RoleID      uint  `json:"role_id" gorm:"not null;index;comment:角色ID"`
}

// 表名
func (AdminUserRole) TableName() string {
	return "admin_user_role"
}
