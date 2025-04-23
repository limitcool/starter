package services

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
	"github.com/spf13/cast"
	"go.uber.org/fx"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo      *repository.RoleRepo
	casbinService casbin.Service
	config        *configs.Config
}

// NewRoleService 创建角色服务
func NewRoleService(params ServiceParams, casbinService casbin.Service) *RoleService {
	// 使用参数中的仓库和配置
	roleRepo := params.RoleRepo
	config := params.Config
	// 获取用户模式
	userMode := enum.GetUserMode(config.Admin.UserMode)

	// 如果是简单模式，返回一个空的实现
	if userMode == enum.UserModeSimple {
		return &RoleService{
			roleRepo:      roleRepo,
			casbinService: nil, // 简单模式不使用Casbin
			config:        config,
		}
	}

	// 分离模式，使用完整的实现
	service := &RoleService{
		roleRepo:      roleRepo,
		casbinService: casbinService,
		config:        config,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "RoleService initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "RoleService stopped")
			return nil
		},
	})

	return service
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, role *model.Role) error {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，直接返回nil
	if userMode == enum.UserModeSimple {
		return nil
	}

	return s.roleRepo.Create(ctx, role)
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, role *model.Role) error {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，直接返回nil
	if userMode == enum.UserModeSimple {
		return nil
	}

	return s.roleRepo.Update(ctx, role)
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, id uint) error {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，直接返回nil
	if userMode == enum.UserModeSimple {
		return nil
	}

	// 检查角色是否已分配给用户
	isAssigned, err := s.roleRepo.IsAssignedToUser(ctx, id)
	if err != nil {
		return errorx.WrapError(err, "检查角色是否已分配给用户失败")
	}
	if isAssigned {
		return errorx.Errorf(errorx.ErrForbidden, "该角色已分配给用户，不能删除")
	}

	// 删除角色菜单关联
	if err := s.roleRepo.DeleteRoleMenus(ctx, id); err != nil {
		return errorx.WrapError(err, "删除角色菜单关联失败")
	}

	// 查询角色信息
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WrapError(err, "查询角色信息失败")
	}

	// 删除Casbin中的角色策略
	if s.casbinService != nil {
		_, err = s.casbinService.DeleteRole(ctx, role.Code)
		if err != nil {
			return errorx.WrapError(err, "删除Casbin角色策略失败")
		}
	}

	// 删除角色
	return s.roleRepo.Delete(ctx, id)
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，返回一个空的角色
	if userMode == enum.UserModeSimple {
		return &model.Role{
			BaseModel: model.BaseModel{ID: id},
			Name:      "管理员",
			Code:      "admin",
		}, nil
	}

	return s.roleRepo.GetByID(ctx, id)
}

// GetRoles 获取角色列表
func (s *RoleService) GetRoles(ctx context.Context) ([]model.Role, error) {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，返回一个包含管理员角色的列表
	if userMode == enum.UserModeSimple {
		return []model.Role{
			{
				BaseModel: model.BaseModel{ID: 1},
				Name:      "管理员",
				Code:      "admin",
			},
		}, nil
	}

	return s.roleRepo.GetAll(ctx)
}

// AssignRolesToUser 为用户分配角色
func (s *RoleService) AssignRolesToUser(ctx context.Context, userID int64, roleIDs []uint) error {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，直接返回nil
	if userMode == enum.UserModeSimple {
		return nil
	}

	// 使用 roleRepo 的 AssignRolesToUser 方法
	if err := s.roleRepo.AssignRolesToUser(ctx, userID, roleIDs); err != nil {
		return errorx.WrapError(err, "为用户分配角色失败")
	}

	// 更新Casbin中的用户角色
	if s.casbinService != nil {
		userIDStr := cast.ToString(userID)

		// 获取用户当前角色
		roles, err := s.casbinService.GetRolesForUser(ctx, userIDStr)
		if err != nil {
			return errorx.WrapError(err, "获取用户角色失败")
		}

		// 移除所有角色
		for _, role := range roles {
			_, err = s.casbinService.DeleteRoleForUser(ctx, userIDStr, role)
			if err != nil {
				return errorx.WrapError(err, "操作失败")
			}
		}

		// 添加新角色
		if len(roleIDs) > 0 {
			for _, roleID := range roleIDs {
				// 查询角色编码
				roleObj, err := s.roleRepo.GetByID(ctx, roleID)
				if err != nil {
					return errorx.WrapError(err, "查询角色信息失败")
				}

				// 添加用户角色关联
				_, err = s.casbinService.AddRoleForUser(ctx, userIDStr, roleObj.Code)
				if err != nil {
					return errorx.WrapError(err, "添加用户角色关联失败")
				}
			}
		}
	}

	return nil
}

// GetUserRoleIDs 获取用户角色ID列表
func (s *RoleService) GetUserRoleIDs(ctx context.Context, userID uint) ([]uint, error) {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，返回一个包含管理员角色ID的列表
	if userMode == enum.UserModeSimple {
		return []uint{1}, nil
	}

	return s.roleRepo.GetRoleIDsByUserID(ctx, userID)
}

// GetRoleMenuIDs 获取角色菜单ID列表
func (s *RoleService) GetRoleMenuIDs(ctx context.Context, roleID uint) ([]uint, error) {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，返回一个空的列表
	if userMode == enum.UserModeSimple {
		return []uint{}, nil
	}

	return s.roleRepo.GetMenuIDsByRoleID(ctx, roleID)
}

// AssignMenusToRole 为角色分配菜单
func (s *RoleService) AssignMenusToRole(ctx context.Context, roleID uint, menuIDs []uint) error {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，直接返回nil
	if userMode == enum.UserModeSimple {
		return nil
	}

	return s.roleRepo.AssignMenusToRole(ctx, roleID, menuIDs)
}

// 为角色设置权限策略
func (s *RoleService) SetRolePermission(ctx context.Context, roleCode string, obj string, act string) error {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，直接返回nil
	if userMode == enum.UserModeSimple || s.casbinService == nil {
		return nil
	}

	_, err := s.casbinService.AddPermissionForRole(ctx, roleCode, obj, act)
	if err != nil {
		return errorx.WrapError(err, "设置角色权限失败")
	}
	return nil
}

// 删除角色的权限策略
func (s *RoleService) DeleteRolePermission(ctx context.Context, roleCode string, obj string, act string) error {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 如果是简单模式，直接返回nil
	if userMode == enum.UserModeSimple || s.casbinService == nil {
		return nil
	}

	_, err := s.casbinService.DeletePermissionForRole(ctx, roleCode, obj, act)
	if err != nil {
		return errorx.WrapError(err, "删除角色权限失败")
	}
	return nil
}
