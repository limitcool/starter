package services

import (
	"context"
	"path/filepath"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// PermissionService 权限服务
type PermissionService struct {
	permissionRepo *repository.PermissionRepo
	roleRepo       *repository.RoleRepo
	menuRepo       *repository.MenuRepo
	casbinService  casbin.Service
	config         *configs.Config
	menuService    MenuServiceInterface
}

// NewPermissionService 创建权限服务
func NewPermissionService(params ServiceParams, casbinService casbin.Service, menuService MenuServiceInterface) *PermissionService {
	// 使用参数中的仓库和配置
	permissionRepo := params.PermissionRepo
	roleRepo := params.RoleRepo
	menuRepo := params.MenuRepo
	config := params.Config
	// 获取用户模式
	userMode := enum.GetUserMode(config.Admin.UserMode)

	// 如果是简单模式，返回一个空的实现
	if userMode == enum.UserModeSimple {
		// 创建 PermissionService 实例
		ps := &PermissionService{
			permissionRepo: permissionRepo,
			roleRepo:       roleRepo,
			menuRepo:       menuRepo,
			casbinService:  nil, // 简单模式不使用Casbin
			config:         config,
			menuService:    menuService,
		}

		return ps
	}

	// 创建 PermissionService 实例
	ps := &PermissionService{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
		menuRepo:       menuRepo,
		casbinService:  casbinService,
		config:         config,
		menuService:    menuService,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "PermissionService initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "PermissionService stopped")
			return nil
		},
	})

	return ps
}

// UpdatePermissionSettings 更新权限系统设置
func (s *PermissionService) UpdatePermissionSettings(ctx context.Context, enabled, defaultAllow bool) error {
	// 更新内存中的配置
	s.config.Casbin.Enabled = enabled
	s.config.Casbin.DefaultAllow = defaultAllow

	// 更新配置文件
	v := viper.New()
	v.SetConfigFile(filepath.Join("configs", "config.yaml"))

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	v.Set("casbin.enabled", enabled)
	v.Set("casbin.default_allow", defaultAllow)

	return v.WriteConfig()
}

// GetPermissions 获取权限列表
func (s *PermissionService) GetPermissions(ctx context.Context) ([]model.Permission, error) {
	return s.permissionRepo.GetAll(ctx)
}

// GetPermissionByID 获取权限详情
func (s *PermissionService) GetPermissionByID(ctx context.Context, id uint) (*model.Permission, error) {
	return s.permissionRepo.GetByID(ctx, id)
}

// GetPermission 获取权限详情
func (s *PermissionService) GetPermission(ctx context.Context, id uint64) (*model.Permission, error) {
	return s.permissionRepo.GetByID(ctx, uint(id))
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(ctx context.Context, permission *model.Permission) error {
	return s.permissionRepo.Create(ctx, permission)
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	return s.permissionRepo.Update(ctx, permission)
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(ctx context.Context, id uint) error {
	return s.permissionRepo.Delete(ctx, id)
}

// AssignPermissionToRole 为角色分配权限
func (s *PermissionService) AssignPermissionToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	// 获取角色
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	// 开始数据库事务
	tx := s.permissionRepo.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务中执行数据库操作
	if err := tx.WithContext(ctx).Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新的角色权限关联
	if len(permissionIDs) > 0 {
		var rolePermissions []model.RolePermission
		for _, permID := range permissionIDs {
			rolePermissions = append(rolePermissions, model.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
			})
		}
		if err := tx.WithContext(ctx).Create(&rolePermissions).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交数据库事务
	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}

	// 使用 Casbin 事务批量更新权限策略
	// 注意：这里我们使用 Casbin 的批量操作 API
	// 1. 删除 Casbin 中的角色权限
	_, err = s.casbinService.DeleteRole(ctx, role.Code)
	if err != nil {
		// 如果 Casbin 操作失败，记录错误但不回滚数据库事务
		// 因为数据库事务已经提交，我们可以在后续的同步操作中修复 Casbin 状态
		logger.ErrorContext(ctx, "删除 Casbin 角色权限失败", "error", err)
		// 这里可以添加重试逻辑或触发异步修复任务
	}

	// 2. 批量添加新的权限策略
	if len(permissionIDs) > 0 {
		// 准备批量添加的策略
		var policies [][]string
		for _, permID := range permissionIDs {
			perm, err := s.permissionRepo.GetByID(ctx, permID)
			if err != nil {
				logger.WarnContext(ctx, "获取权限信息失败", "permission_id", permID, "error", err)
				continue
			}

			// 添加策略 [role.Code, perm.Code, "*"]
			policies = append(policies, []string{role.Code, perm.Code, "*"})
		}

		// 批量添加策略
		if len(policies) > 0 {
			success, err := s.casbinService.AddPolicies(ctx, policies)
			if err != nil || !success {
				logger.ErrorContext(ctx, "批量添加 Casbin 权限策略失败", "error", err)
				// 这里可以添加重试逻辑或触发异步修复任务
			}
		}
	}

	return nil
}

// GetPermissionsByRoleID 获取角色的权限列表
func (s *PermissionService) GetPermissionsByRoleID(ctx context.Context, roleID uint) ([]model.Permission, error) {
	return s.permissionRepo.GetByRoleID(ctx, roleID)
}

// GetPermissionsByUserID 获取用户的权限列表
func (s *PermissionService) GetPermissionsByUserID(ctx context.Context, userID uint) ([]model.Permission, error) {
	return s.permissionRepo.GetByUserID(ctx, userID)
}

// CheckPermission 检查权限
func (s *PermissionService) CheckPermission(ctx context.Context, userID string, obj string, act string) (bool, error) {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，直接返回true
	if userMode == enum.UserModeSimple || s.casbinService == nil {
		return true, nil
	}

	// 分离模式，使用Casbin检查权限
	return s.casbinService.CheckPermission(ctx, userID, obj, act)
}

// AssignRolesToUser 为用户分配角色
func (s *PermissionService) AssignRolesToUser(ctx context.Context, userID string, roleIDs []uint) error {
	// 获取用户当前角色
	currentRoles, err := s.casbinService.GetRolesForUser(ctx, userID)
	if err != nil {
		return err
	}

	// 将 userID 转换为 int64 用于数据库操作
	userIDInt, err := cast.ToInt64E(userID)
	if err != nil {
		return err
	}

	// 开始数据库事务
	tx := s.roleRepo.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务中删除原有的用户角色关联
	if err := tx.WithContext(ctx).Where("user_id = ?", userIDInt).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新的用户角色关联
	if len(roleIDs) > 0 {
		var userRoles []model.UserRole
		for _, roleID := range roleIDs {
			userRoles = append(userRoles, model.UserRole{
				UserID: userIDInt,
				RoleID: roleID,
			})
		}
		if err := tx.WithContext(ctx).Create(&userRoles).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交数据库事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 使用 Casbin 批量操作 API 更新权限
	// 1. 批量删除当前角色
	if len(currentRoles) > 0 {
		success, err := s.casbinService.DeleteRolesForUser(ctx, userID)
		if err != nil || !success {
			logger.ErrorContext(ctx, "批量删除用户角色失败", "error", err)
			// 这里可以添加重试逻辑或触发异步修复任务
		}
	}

	// 2. 批量添加新角色
	if len(roleIDs) > 0 {
		// 准备批量添加的角色
		var roleCodes []string
		for _, roleID := range roleIDs {
			role, err := s.roleRepo.GetByID(ctx, roleID)
			if err != nil {
				logger.WarnContext(ctx, "获取角色信息失败", "role_id", roleID, "error", err)
				continue
			}
			roleCodes = append(roleCodes, role.Code)
		}

		// 批量添加角色
		if len(roleCodes) > 0 {
			success, err := s.casbinService.AddRolesForUser(ctx, userID, roleCodes)
			if err != nil || !success {
				logger.ErrorContext(ctx, "批量添加用户角色失败", "error", err)
				// 这里可以添加重试逻辑或触发异步修复任务
			}
		}
	}

	return nil
}

// GetUserRoles 获取用户角色
func (s *PermissionService) GetUserRoles(ctx context.Context, userID string) ([]model.Role, error) {
	// 获取用户角色编码
	roleCodes, err := s.casbinService.GetRolesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取角色详情
	var roles []model.Role
	for _, code := range roleCodes {
		role, err := s.roleRepo.GetByCode(ctx, code)
		if err != nil {
			continue
		}
		roles = append(roles, *role)
	}

	return roles, nil
}

// GetUserMenus 获取用户菜单
func (s *PermissionService) GetUserMenus(ctx context.Context, userID int64) ([]*model.MenuTree, error) {
	// 将int64转换为string
	userIDStr := cast.ToString(userID)

	// 获取用户角色
	roles, err := s.GetUserRoles(ctx, userIDStr)
	if err != nil {
		return nil, err
	}

	// 使用 MenuService 的方法获取用户菜单树
	return s.menuService.GetUserMenuTree(ctx, userIDStr, roles)
}
