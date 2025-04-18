package repository

import (
	"github.com/limitcool/starter/internal/model"
)

// MenuRepository 菜单仓库接口
type MenuRepository interface {
	GetByID(id uint) (*model.Menu, error)
	GetAll() ([]*model.Menu, error)
	Create(menu *model.Menu) error
	Update(menu *model.Menu) error
	Delete(id uint) error
	GetByRoleID(roleID uint) ([]*model.Menu, error)
	GetByUserID(userID uint) ([]*model.Menu, error)
	GetPermsByUserID(userID uint) ([]string, error)
	GetPermsByRoleIDs(roleIDs []uint) ([]string, error)
	AssignMenuToRole(roleID uint, menuIDs []uint) error
	BuildMenuTree(menus []*model.Menu) []*model.MenuTree
}

// RoleRepository 角色仓库接口
type RoleRepository interface {
	GetByID(id uint) (*model.Role, error)
	GetAll() ([]model.Role, error)
	Create(role *model.Role) error
	Update(role *model.Role) error
	Delete(id uint) error
	IsAssignedToUser(id uint) (bool, error)
	DeleteRoleMenus(roleID uint) error
	GetMenuIDsByRoleID(roleID uint) ([]uint, error)
	GetRoleIDsByUserID(userID uint) ([]uint, error)
	AssignRolesToUser(userID int64, roleIDs []uint) error
}

// SysUserRepository 系统用户仓库接口
type SysUserRepository interface {
	GetByID(id int64) (*model.SysUser, error)
	GetByUsername(username string) (*model.SysUser, error)
	Create(user *model.SysUser) error
	Update(user *model.SysUser) error
	UpdateFields(id int64, fields map[string]interface{}) error
	Delete(id int64) error
	List(query *model.SysUserQuery) ([]model.SysUser, int64, error)
	UpdateAvatar(userID int64, fileID uint) error
}

// PermissionRepository 权限仓库接口
type PermissionRepository interface {
	GetByID(id uint) (*model.Permission, error)
	GetAll() ([]model.Permission, error)
	Create(permission *model.Permission) error
	Update(permission *model.Permission) error
	Delete(id uint) error
	GetByRoleID(roleID uint) ([]model.Permission, error)
	GetByUserID(userID uint) ([]model.Permission, error)
	AssignPermissionToRole(roleID uint, permissionIDs []uint) error
}

// OperationLogRepository 操作日志仓库接口
type OperationLogRepository interface {
	Create(log *model.OperationLog) error
	GetLogs(query *OperationLogQuery) (*PageResult, error)
	Delete(id uint) error
	BatchDelete(ids []uint) error
}

// FileRepository 文件仓库接口
type FileRepository interface {
	GetByID(id string) (*model.File, error)
	Create(file *model.File) error
	Update(file *model.File) error
	Delete(id string) error
	UpdateUserAvatar(userID int64, fileID uint) error
	UpdateSysUserAvatar(userID int64, fileID uint) error
}
