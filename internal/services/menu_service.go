package services

import (
	"context"
	"strconv"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
	"go.uber.org/fx"
)

// MenuService 菜单服务
type MenuService struct {
	menuRepo      *repository.MenuRepo
	roleRepo      *repository.RoleRepo
	casbinService casbin.Service
}

// NewMenuService 创建菜单服务
func NewMenuService(params ServiceParams, casbinService casbin.Service) *MenuService {
	// 使用参数中的仓库和配置
	menuRepo := params.MenuRepo
	roleRepo := params.RoleRepo
	config := params.Config
	// 获取用户模式
	userMode := enum.GetUserMode(config.Admin.UserMode)

	// 如果是简单模式，返回一个空的实现
	if userMode == enum.UserModeSimple {
		return &MenuService{
			menuRepo:      menuRepo,
			roleRepo:      roleRepo,
			casbinService: nil, // 简单模式不使用Casbin
		}
	}

	// 分离模式，使用完整的实现
	service := &MenuService{
		menuRepo:      menuRepo,
		roleRepo:      roleRepo,
		casbinService: casbinService,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "MenuService initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "MenuService stopped")
			return nil
		},
	})

	return service
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(ctx context.Context, menu *model.Menu) error {
	err := s.menuRepo.Create(ctx, menu)
	if err != nil {
		return errorx.WrapError(err, "创建菜单失败")
	}
	return nil
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(ctx context.Context, menu *model.Menu) error {
	err := s.menuRepo.Update(ctx, menu)
	if err != nil {
		return errorx.WrapError(err, "更新菜单失败")
	}
	return nil
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(ctx context.Context, id uint) error {
	err := s.menuRepo.Delete(ctx, id)
	if err != nil {
		return errorx.WrapError(err, "删除菜单失败")
	}
	return nil
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(ctx context.Context, id uint) (*model.Menu, error) {
	menu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.WrapError(err, "获取菜单失败")
	}
	return menu, nil
}

// GetAllMenus 获取所有菜单
func (s *MenuService) GetAllMenus(ctx context.Context) ([]*model.Menu, error) {
	menus, err := s.menuRepo.GetAll(ctx)
	if err != nil {
		return nil, errorx.WrapError(err, "获取所有菜单失败")
	}
	return menus, nil
}

// GetMenusByRoleID 获取角色菜单
func (s *MenuService) GetMenusByRoleID(ctx context.Context, roleID uint) ([]*model.Menu, error) {
	menus, err := s.menuRepo.GetByRoleID(ctx, roleID)
	if err != nil {
		return nil, errorx.WrapError(err, "获取角色菜单失败")
	}
	return menus, nil
}

// GetUserMenus 获取用户菜单
func (s *MenuService) GetUserMenus(ctx context.Context, userID uint) ([]*model.Menu, error) {
	menus, err := s.menuRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapError(err, "获取用户菜单失败")
	}
	return menus, nil
}

// AssignMenuToRole 为角色分配菜单
func (s *MenuService) AssignMenuToRole(ctx context.Context, roleID uint, menuIDs []uint) error {
	err := s.menuRepo.AssignMenuToRole(ctx, roleID, menuIDs)
	if err != nil {
		return errorx.WrapError(err, "为角色分配菜单失败")
	}

	// 更新 Casbin 策略
	// 获取角色的所有权限
	perms, err := s.menuRepo.GetPermsByRoleIDs(ctx, []uint{roleID})
	if err != nil {
		return errorx.WrapError(err, "获取角色权限失败")
	}

	// 更新 Casbin 策略
	// 先获取角色编码
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return errorx.WrapError(err, "获取角色信息失败")
	}

	// 删除原有策略
	policies, err := s.casbinService.GetPermissionsForRole(ctx, role.Code)
	if err != nil {
		return errorx.WrapError(err, "获取角色权限策略失败")
	}

	// 如果有原有策略，删除它们
	if len(policies) > 0 {
		_, err = s.casbinService.RemovePolicies(ctx, policies)
		if err != nil {
			return errorx.WrapError(err, "删除原有权限策略失败")
		}
	}

	// 添加新策略
	if len(perms) > 0 {
		var newPolicies [][]string
		for _, perm := range perms {
			newPolicies = append(newPolicies, []string{role.Code, perm, "*"})
		}

		_, err = s.casbinService.AddPolicies(ctx, newPolicies)
		if err != nil {
			return errorx.WrapError(err, "添加权限策略失败")
		}
	}

	return nil
}

// GetMenuTree 获取菜单树(用于前端菜单选择)
func (s *MenuService) GetMenuTree(ctx context.Context) ([]*model.MenuTree, error) {
	menus, err := s.menuRepo.GetAll(ctx)
	if err != nil {
		return nil, errorx.WrapError(err, "获取所有菜单失败")
	}
	return s.menuRepo.BuildMenuTree(ctx, menus), nil
}

// GetUserMenuTree 获取用户菜单树
func (s *MenuService) GetUserMenuTree(ctx context.Context, userID string, roles []model.Role) ([]*model.MenuTree, error) {
	// 获取角色菜单
	var menuIDs []uint
	for _, role := range roles {
		// 管理员角色获取所有菜单
		if role.Code == "admin" {
			allMenus, err := s.menuRepo.GetAll(ctx)
			if err != nil {
				return nil, errorx.WrapError(err, "获取所有菜单失败")
			}
			return s.menuRepo.BuildMenuTree(ctx, allMenus), nil
		}

		// 获取角色菜单
		roleMenus, err := s.menuRepo.GetByRoleID(ctx, role.ID)
		if err != nil {
			continue
		}

		for _, menu := range roleMenus {
			menuIDs = append(menuIDs, menu.ID)
		}
	}

	// 去重
	uniqueMenuIDs := make(map[uint]bool)
	var uniqueIDs []uint
	for _, id := range menuIDs {
		if !uniqueMenuIDs[id] {
			uniqueMenuIDs[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	// 获取菜单详情
	var menus []*model.Menu
	for _, id := range uniqueIDs {
		menu, err := s.menuRepo.GetByIDWithRelations(ctx, id)
		if err != nil {
			continue
		}
		menus = append(menus, menu)
	}

	// 构建菜单树
	return s.menuRepo.BuildMenuTree(ctx, menus), nil
}

// GetMenuPermsByUserID 获取用户菜单权限标识
func (s *MenuService) GetMenuPermsByUserID(ctx context.Context, userID uint) ([]string, error) {
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 获取用户角色
	roles, err := s.casbinService.GetRolesForUser(ctx, userIDStr)
	if err != nil {
		return nil, errorx.WrapError(err, "获取用户角色失败")
	}

	if len(roles) == 0 {
		return []string{}, nil
	}

	// 提取角色ID
	var roleIDs []uint
	for _, role := range roles {
		id, _ := strconv.ParseUint(role, 10, 64)
		roleIDs = append(roleIDs, uint(id))
	}

	// 查询角色菜单权限标识
	return s.menuRepo.GetPermsByRoleIDs(ctx, roleIDs)
}
