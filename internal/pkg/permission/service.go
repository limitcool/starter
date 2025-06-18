package permission

import (
	"context"
	"strings"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// Service 权限服务
type Service struct {
	casbinService *casbin.Service
	userRepo      *model.UserRepo
	roleRepo      *model.RoleRepo
	permRepo      *model.PermissionRepo
	menuRepo      *model.MenuRepo
}

// NewService 创建权限服务
func NewService(
	casbinService *casbin.Service,
	userRepo *model.UserRepo,
	roleRepo *model.RoleRepo,
	permRepo *model.PermissionRepo,
	menuRepo *model.MenuRepo,
) *Service {
	return &Service{
		casbinService: casbinService,
		userRepo:      userRepo,
		roleRepo:      roleRepo,
		permRepo:      permRepo,
		menuRepo:      menuRepo,
	}
}

// CheckPermission 检查用户权限
func (s *Service) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	// 检查用户是否是管理员
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, errorx.WrapError(err, "获取用户信息失败")
	}

	// 管理员拥有所有权限
	if user.IsAdmin {
		return true, nil
	}

	// 如果Casbin未启用，使用简单的角色权限检查
	if s.casbinService == nil {
		return s.checkPermissionByRole(ctx, userID, resource, action)
	}

	// 使用Casbin进行权限验证
	userKey := casbin.GetUserKey(userID)
	resourceKey := casbin.GetResourceKey(resource)

	return s.casbinService.Enforce(ctx, userKey, resourceKey, action)
}

// checkPermissionByRole 通过角色检查权限（当Casbin未启用时）
func (s *Service) checkPermissionByRole(ctx context.Context, userID int64, resource, action string) (bool, error) {
	// 获取用户角色
	roles, err := s.roleRepo.GetRolesByUserID(ctx, userID)
	if err != nil {
		return false, errorx.WrapError(err, "获取用户角色失败")
	}

	// 检查角色权限
	for _, role := range roles {
		permissions, err := s.roleRepo.GetPermissionsByRoleID(ctx, role.ID)
		if err != nil {
			continue
		}

		for _, permission := range permissions {
			if permission.Key == resource+":"+action {
				return true, nil
			}
		}
	}

	return false, nil
}

// AssignRolesToUser 为用户分配角色
func (s *Service) AssignRolesToUser(ctx context.Context, userID int64, roleKeys []string) error {
	if s.casbinService == nil {
		return errorx.WrapError(nil, "Casbin服务未启用")
	}

	userKey := casbin.GetUserKey(userID)

	// 删除用户的所有角色
	err := s.casbinService.DeleteRolesForUser(ctx, userKey)
	if err != nil {
		return errorx.WrapError(err, "删除用户角色失败")
	}

	// 添加新角色
	for _, roleKey := range roleKeys {
		err = s.casbinService.AddRoleForUser(ctx, userKey, roleKey)
		if err != nil {
			logger.ErrorContext(ctx, "添加用户角色失败", "user_id", userID, "role_key", roleKey, "error", err)
		}
	}

	return nil
}

// syncUserRolesToCasbin 同步用户角色到Casbin
func (s *Service) syncUserRolesToCasbin(ctx context.Context, userID int64, roleIDs []uint) error {
	userKey := casbin.GetUserKey(userID)

	// 删除用户的所有角色
	err := s.casbinService.DeleteRolesForUser(ctx, userKey)
	if err != nil {
		return errorx.WrapError(err, "清除用户角色失败")
	}

	// 添加新角色
	for _, roleID := range roleIDs {
		roleKey := casbin.GetRoleKey(roleID)
		err = s.casbinService.AddRoleForUser(ctx, userKey, roleKey)
		if err != nil {
			return errorx.WrapError(err, "添加用户角色失败")
		}
	}

	return nil
}

// AssignPermissionsToRole 为角色分配权限
func (s *Service) AssignPermissionsToRole(ctx context.Context, roleID uint, permissionKeys []string) error {
	if s.casbinService == nil {
		return errorx.WrapError(nil, "Casbin服务未启用")
	}

	roleKey := casbin.GetRoleKey(roleID)

	// 删除角色的所有权限策略
	// TODO: 需要实现删除角色所有权限的方法

	// 添加新权限策略
	for _, permissionKey := range permissionKeys {
		// 解析权限key，格式如 "user:create"
		parts := strings.Split(permissionKey, ":")
		if len(parts) != 2 {
			continue
		}
		resource := parts[0]
		action := parts[1]

		err := s.casbinService.AddPolicy(ctx, roleKey, resource, action)
		if err != nil {
			logger.ErrorContext(ctx, "添加角色权限策略失败", "role_id", roleID, "permission_key", permissionKey, "error", err)
		}
	}

	return nil
}

// 注意：在新的设计中，不再需要这个方法，权限直接通过Casbin管理

// GetUserRoles 获取用户角色列表
func (s *Service) GetUserRoles(ctx context.Context, userID int64) ([]model.Role, error) {
	if s.casbinService == nil {
		return nil, errorx.WrapError(nil, "Casbin服务未启用")
	}

	userKey := casbin.GetUserKey(userID)
	roleKeys, err := s.casbinService.GetRolesForUser(ctx, userKey)
	if err != nil {
		return nil, errorx.WrapError(err, "获取用户角色失败")
	}

	// 根据角色Key查询角色详情
	var roles []model.Role
	for _, roleKey := range roleKeys {
		role, err := s.roleRepo.GetByKey(ctx, roleKey)
		if err == nil {
			roles = append(roles, *role)
		}
	}

	return roles, nil
}

// GetUserPermissions 获取用户权限列表
func (s *Service) GetUserPermissions(ctx context.Context, userID int64) ([]model.Permission, error) {
	if s.casbinService == nil {
		return nil, errorx.WrapError(nil, "Casbin服务未启用")
	}

	userKey := casbin.GetUserKey(userID)
	permissionPolicies, err := s.casbinService.GetImplicitPermissionsForUser(ctx, userKey)
	if err != nil {
		return nil, errorx.WrapError(err, "获取用户权限失败")
	}

	// 根据权限策略查询权限详情
	var permissions []model.Permission
	for _, policy := range permissionPolicies {
		if len(policy) >= 2 {
			// policy格式: [subject, object, action]
			permissionKey := policy[1] + ":" + policy[2]
			permission, err := s.permRepo.GetByKey(ctx, permissionKey)
			if err == nil {
				permissions = append(permissions, *permission)
			}
		}
	}

	return permissions, nil
}

// GetUserMenus 获取用户可访问的菜单
func (s *Service) GetUserMenus(ctx context.Context, userID int64) ([]model.Menu, error) {
	return s.GetUserMenusByPlatform(ctx, userID, "admin")
}

// GetUserMenusByPlatform 根据平台获取用户可访问的菜单
func (s *Service) GetUserMenusByPlatform(ctx context.Context, userID int64, platform string) ([]model.Menu, error) {
	// 检查用户是否是管理员
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapError(err, "获取用户信息失败")
	}

	// 获取指定平台的所有可见菜单
	opts := &model.QueryOptions{
		Condition: "is_visible = ? AND platform = ?",
		Args:      []any{true, platform},
	}
	allMenus, err := s.menuRepo.GetAll(ctx, opts)
	if err != nil {
		return nil, errorx.WrapError(err, "获取菜单失败")
	}

	var accessibleMenus []model.Menu

	// 管理员可以访问所有菜单
	if user.IsAdmin {
		accessibleMenus = allMenus
	} else {
		// 普通用户根据权限过滤菜单
		for _, menu := range allMenus {
			if menu.PermissionKey == "" {
				// 没有权限要求的菜单，所有人都可以访问
				accessibleMenus = append(accessibleMenus, menu)
				continue
			}

			// 检查用户是否有该菜单要求的权限
			parts := strings.Split(menu.PermissionKey, ":")
			if len(parts) == 2 {
				resource := parts[0]
				action := parts[1]
				hasPermission, err := s.CheckPermission(ctx, userID, resource, action)
				if err == nil && hasPermission {
					accessibleMenus = append(accessibleMenus, menu)
				}
			}
		}
	}

	// 构建菜单树
	return s.menuRepo.BuildMenuTree(accessibleMenus, 0), nil
}

// IsAdmin 检查用户是否是管理员
func (s *Service) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, errorx.WrapError(err, "获取用户信息失败")
	}
	return user.IsAdmin, nil
}

// HasPermission 检查用户是否有特定权限（别名方法）
func (s *Service) HasPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	return s.CheckPermission(ctx, userID, resource, action)
}

// InitializePermissions 初始化权限系统
func (s *Service) InitializePermissions(ctx context.Context) error {
	if s.casbinService == nil {
		logger.InfoContext(ctx, "Casbin未启用，跳过权限初始化")
		return nil
	}

	// 同步所有角色权限到Casbin
	roles, err := s.roleRepo.GetEnabledRoles(ctx)
	if err != nil {
		return errorx.WrapError(err, "获取角色列表失败")
	}

	for _, role := range roles {
		permissions, err := s.roleRepo.GetPermissionsByRoleID(ctx, role.ID)
		if err != nil {
			logger.ErrorContext(ctx, "获取角色权限失败", "role_id", role.ID, "error", err)
			continue
		}

		roleKey := casbin.GetRoleKey(role.ID)
		for _, permission := range permissions {
			// 解析权限key，格式如 "user:create"
			parts := strings.Split(permission.Key, ":")
			if len(parts) != 2 {
				continue
			}
			resource := parts[0]
			action := parts[1]

			err = s.casbinService.AddPolicy(ctx, roleKey, resource, action)
			if err != nil {
				logger.ErrorContext(ctx, "添加角色权限策略失败", "role_id", role.ID, "permission_id", permission.ID, "error", err)
			}
		}
	}

	logger.InfoContext(ctx, "权限系统初始化完成", "roles_count", len(roles))
	return nil
}
