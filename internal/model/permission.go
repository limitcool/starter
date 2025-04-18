package model

import (
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/storage/sqldb"
)

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

// Create 创建权限
func (p *Permission) Create() error {
	return sqldb.Instance().DB().Create(p).Error
}

// Update 更新权限
func (p *Permission) Update() error {
	return sqldb.Instance().DB().Model(&Permission{}).Where("id = ?", p.ID).Updates(p).Error
}

// Delete 删除权限
func (p *Permission) Delete(id uint) error {
	return sqldb.Instance().DB().Delete(&Permission{}, id).Error
}

// GetByID 根据ID获取权限
func (p *Permission) GetByID(id uint) (*Permission, error) {
	var permission Permission
	err := sqldb.Instance().DB().Where("id = ?", id).First(&permission).Error
	return &permission, err
}

// GetAll 获取所有权限
func (p *Permission) GetAll() ([]Permission, error) {
	var permissions []Permission
	err := sqldb.Instance().DB().Find(&permissions).Error
	return permissions, err
}

// GetByRoleID 获取角色的权限列表
func (p *Permission) GetByRoleID(roleID uint) ([]Permission, error) {
	var permissions []Permission
	db := sqldb.Instance().DB()

	// 通过关联表查询
	err := db.Joins("JOIN sys_role_permission ON sys_role_permission.permission_id = sys_permission.id").
		Where("sys_role_permission.role_id = ?", roleID).
		Find(&permissions).Error

	return permissions, err
}

// GetByUserID 获取用户的权限列表
func (p *Permission) GetByUserID(userID uint) ([]Permission, error) {
	var permissions []Permission
	db := sqldb.Instance().DB()

	// 通过用户角色关联查询权限
	err := db.Joins("JOIN sys_role_permission ON sys_role_permission.permission_id = sys_permission.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_permission.role_id").
		Where("sys_user_role.user_id = ?", userID).
		Find(&permissions).Error

	return permissions, err
}

// BatchCreate 批量创建角色权限关联
func (rp *RolePermission) BatchCreate(rolePermissions []RolePermission) error {
	return sqldb.Instance().DB().Create(&rolePermissions).Error
}

// DeleteByRoleID 删除角色的权限关联
func (rp *RolePermission) DeleteByRoleID(roleID uint) error {
	return sqldb.Instance().DB().Where("role_id = ?", roleID).Delete(&RolePermission{}).Error
}
