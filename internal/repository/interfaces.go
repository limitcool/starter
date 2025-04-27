package repository

import (
	"context"

	"github.com/limitcool/starter/internal/model"
)

// MenuRepository 菜单仓库接口
type MenuRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Menu, error)
	GetAll(ctx context.Context) ([]*model.Menu, error)
	Create(ctx context.Context, menu *model.Menu) error
	Update(ctx context.Context, menu *model.Menu) error
	Delete(ctx context.Context, id uint) error
	GetByRoleID(ctx context.Context, roleID uint) ([]*model.Menu, error)
	GetByUserID(ctx context.Context, userID uint) ([]*model.Menu, error)
	GetPermsByUserID(ctx context.Context, userID uint) ([]string, error)
	GetPermsByRoleIDs(ctx context.Context, roleIDs []uint) ([]string, error)
	AssignMenuToRole(ctx context.Context, roleID uint, menuIDs []uint) error
	BuildMenuTree(menus []*model.Menu) []*model.MenuTree // 这个方法不需要上下文，因为它只是一个工具方法
}

// RoleRepository 角色仓库接口
type RoleRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Role, error)
	GetAll(ctx context.Context) ([]model.Role, error)
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id uint) error
	IsAssignedToUser(ctx context.Context, id uint) (bool, error)
	DeleteRoleMenus(ctx context.Context, roleID uint) error
	GetMenuIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error)
	GetRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error)
	AssignRolesToUser(ctx context.Context, userID int64, roleIDs []uint) error
}

// AdminUserRepository 管理员用户仓库接口
type AdminUserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.AdminUser, error)
	GetByUsername(ctx context.Context, username string) (*model.AdminUser, error)
	Create(ctx context.Context, user *model.AdminUser) error
	Update(ctx context.Context, user *model.AdminUser) error
	UpdateFields(ctx context.Context, id int64, fields map[string]any) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, query *model.AdminUserQuery) ([]model.AdminUser, int64, error)
	UpdateAvatar(ctx context.Context, userID int64, fileID uint) error
}

// PermissionRepository 权限仓库接口
type PermissionRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Permission, error)
	GetAll(ctx context.Context) ([]model.Permission, error)
	Create(ctx context.Context, permission *model.Permission) error
	Update(ctx context.Context, permission *model.Permission) error
	Delete(ctx context.Context, id uint) error
	GetByRoleID(ctx context.Context, roleID uint) ([]model.Permission, error)
	GetByUserID(ctx context.Context, userID uint) ([]model.Permission, error)
	AssignPermissionToRole(ctx context.Context, roleID uint, permissionIDs []uint) error
}

// OperationLogRepository 操作日志仓库接口
type OperationLogRepository interface {
	Create(ctx context.Context, log *model.OperationLog) error
	GetLogs(ctx context.Context, query *OperationLogQuery) (*PageResult, error)
	Delete(ctx context.Context, id uint) error
	BatchDelete(ctx context.Context, ids []uint) error
}

// FileRepository 文件仓库接口
type FileRepository interface {
	GetByID(ctx context.Context, id string) (*model.File, error)
	Create(ctx context.Context, file *model.File) error
	Update(ctx context.Context, file *model.File) error
	Delete(ctx context.Context, id string) error
	UpdateUserAvatar(ctx context.Context, userID int64, fileID uint) error
	UpdateAdminUserAvatar(ctx context.Context, userID int64, fileID uint) error
}
