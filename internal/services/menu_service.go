package services

import (
	"strconv"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/repository"
)

// MenuService 菜单服务
type MenuService struct {
	menuRepo      *repository.MenuRepo
	roleRepo      *repository.RoleRepo
	casbinService casbin.Service
}

// NewMenuService 创建菜单服务
func NewMenuService(menuRepo *repository.MenuRepo, roleRepo *repository.RoleRepo, casbinService casbin.Service) *MenuService {
	return &MenuService{
		menuRepo:      menuRepo,
		roleRepo:      roleRepo,
		casbinService: casbinService,
	}
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(menu *model.Menu) error {
	err := s.menuRepo.Create(menu)
	if err != nil {
		return errorx.WrapError(err, "创建菜单失败")
	}
	return nil
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(menu *model.Menu) error {
	err := s.menuRepo.Update(menu)
	if err != nil {
		return errorx.WrapError(err, "更新菜单失败")
	}
	return nil
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	err := s.menuRepo.Delete(id)
	if err != nil {
		return errorx.WrapError(err, "删除菜单失败")
	}
	return nil
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(id uint) (*model.Menu, error) {
	menu, err := s.menuRepo.GetByID(id)
	if err != nil {
		return nil, errorx.WrapError(err, "获取菜单失败")
	}
	return menu, nil
}

// GetAllMenus 获取所有菜单
func (s *MenuService) GetAllMenus() ([]*model.Menu, error) {
	menus, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, errorx.WrapError(err, "获取所有菜单失败")
	}
	return menus, nil
}

// GetMenusByRoleID 获取角色菜单
func (s *MenuService) GetMenusByRoleID(roleID uint) ([]*model.Menu, error) {
	menus, err := s.menuRepo.GetByRoleID(roleID)
	if err != nil {
		return nil, errorx.WrapError(err, "获取角色菜单失败")
	}
	return menus, nil
}

// GetUserMenus 获取用户菜单
func (s *MenuService) GetUserMenus(userID uint) ([]*model.Menu, error) {
	menus, err := s.menuRepo.GetByUserID(userID)
	if err != nil {
		return nil, errorx.WrapError(err, "获取用户菜单失败")
	}
	return menus, nil
}

// AssignMenuToRole 为角色分配菜单
func (s *MenuService) AssignMenuToRole(roleID uint, menuIDs []uint) error {
	err := s.menuRepo.AssignMenuToRole(roleID, menuIDs)
	if err != nil {
		return errorx.WrapError(err, "为角色分配菜单失败")
	}

	// 更新 Casbin 策略
	// 获取角色的所有权限
	perms, err := s.menuRepo.GetPermsByRoleIDs([]uint{roleID})
	if err != nil {
		return errorx.WrapError(err, "获取角色权限失败")
	}

	// 更新 Casbin 策略
	// 先获取角色编码
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return errorx.WrapError(err, "获取角色信息失败")
	}

	// 删除原有策略
	policies, err := s.casbinService.GetPermissionsForRole(role.Code)
	if err != nil {
		return errorx.WrapError(err, "获取角色权限策略失败")
	}

	// 如果有原有策略，删除它们
	if len(policies) > 0 {
		_, err = s.casbinService.RemovePolicies(policies)
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

		_, err = s.casbinService.AddPolicies(newPolicies)
		if err != nil {
			return errorx.WrapError(err, "添加权限策略失败")
		}
	}

	return nil
}

// GetMenuTree 获取菜单树(用于前端菜单选择)
func (s *MenuService) GetMenuTree() ([]*model.MenuTree, error) {
	menus, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, errorx.WrapError(err, "获取所有菜单失败")
	}
	return s.menuRepo.BuildMenuTree(menus), nil
}

// GetMenuPermsByUserID 获取用户菜单权限标识
func (s *MenuService) GetMenuPermsByUserID(userID uint) ([]string, error) {
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 获取用户角色
	roles, err := s.casbinService.GetRolesForUser(userIDStr)
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
	return s.menuRepo.GetPermsByRoleIDs(roleIDs)
}
