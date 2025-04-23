package services

import (
	"context"

	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/model"
)

// UserServiceInterface 用户服务接口
// 定义用户服务的通用方法，不同模式下的实现可能不同
type UserServiceInterface interface {
	// 用户相关方法
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	Register(ctx context.Context, req v1.UserRegisterRequest, registerIP string, isAdmin bool) (*model.User, error)
	Login(ctx context.Context, username, password string, ip string) (*v1.LoginResponse, error)
	UpdateUser(ctx context.Context, id uint, data map[string]any) error
	ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error
}

// AdminUserServiceInterface 管理员用户服务接口
// 定义管理员用户服务的通用方法，不同模式下的实现可能不同
type AdminUserServiceInterface interface {
	// 管理员登录
	Login(ctx context.Context, username, password string, ip string) (*v1.LoginResponse, error)
	// 刷新令牌
	RefreshToken(ctx context.Context, refreshToken string) (*v1.LoginResponse, error)
	// 获取管理员用户信息
	GetUserInfo(ctx context.Context, id int64) (interface{}, error)
}

// RoleServiceInterface 角色服务接口
type RoleServiceInterface interface {
	// 角色相关方法
	CreateRole(ctx context.Context, role *model.Role) error
	UpdateRole(ctx context.Context, role *model.Role) error
	DeleteRole(ctx context.Context, id uint) error
	GetRoleByID(ctx context.Context, id uint) (*model.Role, error)
	GetRoles(ctx context.Context) ([]model.Role, error)
	AssignRolesToUser(ctx context.Context, userID int64, roleIDs []uint) error
	GetUserRoleIDs(ctx context.Context, userID uint) ([]uint, error)
	GetRoleMenuIDs(ctx context.Context, roleID uint) ([]uint, error)
	AssignMenusToRole(ctx context.Context, roleID uint, menuIDs []uint) error
	// 权限相关方法
	SetRolePermission(ctx context.Context, roleCode, obj, act string) error
	DeleteRolePermission(ctx context.Context, roleCode, obj, act string) error
}

// MenuServiceInterface 菜单服务接口
type MenuServiceInterface interface {
	// 菜单相关方法
	CreateMenu(ctx context.Context, menu *model.Menu) error
	UpdateMenu(ctx context.Context, menu *model.Menu) error
	DeleteMenu(ctx context.Context, id uint) error
	GetMenuByID(ctx context.Context, id uint) (*model.Menu, error)
	GetMenus(ctx context.Context) ([]model.Menu, error)
	GetMenuTree(ctx context.Context) ([]*model.MenuTree, error)
	GetUserMenus(ctx context.Context, userID int64) ([]*model.MenuTree, error)
	// 角色菜单相关方法
	GetMenusByRoleID(ctx context.Context, roleID uint) ([]model.Menu, error)
	AssignMenuToRole(ctx context.Context, roleID uint, menuIDs []uint) error
	// 权限相关方法
	GetMenuPermsByUserID(ctx context.Context, userID uint) ([]string, error)
	// 获取用户菜单树
	GetUserMenuTree(ctx context.Context, userID string, roles []model.Role) ([]*model.MenuTree, error)
}

// PermissionServiceInterface 权限服务接口
type PermissionServiceInterface interface {
	// 权限相关方法
	CreatePermission(ctx context.Context, permission *model.Permission) error
	UpdatePermission(ctx context.Context, permission *model.Permission) error
	DeletePermission(ctx context.Context, id uint) error
	GetPermissionByID(ctx context.Context, id uint) (*model.Permission, error)
	GetPermissions(ctx context.Context) ([]model.Permission, error)
	CheckPermission(ctx context.Context, userID string, obj string, act string) (bool, error)
	GetUserMenus(ctx context.Context, userID int64) ([]*model.MenuTree, error)

	// 扩展方法
	UpdatePermissionSettings(ctx context.Context, enabled bool, defaultAllow bool) error
	GetPermission(ctx context.Context, id uint64) (*model.Permission, error)
	GetPermissionsByUserID(ctx context.Context, userID uint) ([]model.Permission, error)
	GetUserRoles(ctx context.Context, userID string) ([]model.Role, error)
	AssignRolesToUser(ctx context.Context, userID string, roleIDs []uint) error
	AssignPermissionToRole(ctx context.Context, roleID uint, permissionIDs []uint) error
	GetPermissionsByRoleID(ctx context.Context, roleID uint) ([]model.Permission, error)
}
