package services

import (
	"strconv"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/repository"
)

// MenuService 菜单服务
type MenuService struct {
	menuRepo      repository.MenuRepository
	casbinService *CasbinService
}

// NewMenuService 创建菜单服务
func NewMenuService(menuRepo repository.MenuRepository, casbinService *CasbinService) *MenuService {
	return &MenuService{
		menuRepo:      menuRepo,
		casbinService: casbinService,
	}
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(menu *model.Menu) error {
	return s.menuRepo.Create(menu)
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(menu *model.Menu) error {
	return s.menuRepo.Update(menu)
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	return s.menuRepo.Delete(id)
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(id uint) (*model.Menu, error) {
	return s.menuRepo.GetByID(id)
}

// GetAllMenus 获取所有菜单
func (s *MenuService) GetAllMenus() ([]*model.Menu, error) {
	return s.menuRepo.GetAll()
}

// GetMenusByRoleID 获取角色菜单
func (s *MenuService) GetMenusByRoleID(roleID uint) ([]*model.Menu, error) {
	return s.menuRepo.GetByRoleID(roleID)
}

// GetUserMenus 获取用户菜单
func (s *MenuService) GetUserMenus(userID uint) ([]*model.Menu, error) {
	return s.menuRepo.GetByUserID(userID)
}

// AssignMenuToRole 为角色分配菜单
func (s *MenuService) AssignMenuToRole(roleID uint, menuIDs []uint) error {
	return s.menuRepo.AssignMenuToRole(roleID, menuIDs)
}

// GetMenuTree 获取菜单树(用于前端菜单选择)
func (s *MenuService) GetMenuTree() ([]*model.MenuTree, error) {
	menus, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return s.menuRepo.BuildMenuTree(menus), nil
}

// GetMenuPermsByUserID 获取用户菜单权限标识
func (s *MenuService) GetMenuPermsByUserID(userID uint) ([]string, error) {
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 获取用户角色
	roles, err := s.casbinService.GetRolesForUser(userIDStr)
	if err != nil {
		return nil, err
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
