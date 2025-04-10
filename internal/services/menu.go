package services

import (
	"strconv"

	"github.com/limitcool/starter/internal/model"
)

// MenuService 菜单服务
type MenuService struct {
}

// NewMenuService 创建菜单服务
func NewMenuService() *MenuService {
	return &MenuService{}
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(menu *model.Menu) error {
	return menu.Create()
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(menu *model.Menu) error {
	return menu.Update()
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	menu := &model.Menu{}
	return menu.Delete(id)
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(id uint) (*model.Menu, error) {
	menu := &model.Menu{}
	return menu.GetByID(id)
}

// GetAllMenus 获取所有菜单
func (s *MenuService) GetAllMenus() ([]*model.Menu, error) {
	menu := &model.Menu{}
	return menu.GetAll()
}

// GetMenusByRoleID 获取角色菜单
func (s *MenuService) GetMenusByRoleID(roleID uint) ([]*model.Menu, error) {
	menu := &model.Menu{}
	return menu.GetByRoleID(roleID)
}

// GetUserMenus 获取用户菜单
func (s *MenuService) GetUserMenus(userID uint) ([]*model.Menu, error) {
	menu := &model.Menu{}
	return menu.GetByUserID(userID)
}

// 为角色分配菜单
func (s *MenuService) AssignMenuToRole(roleID uint, menuIDs []uint) error {
	roleMenu := &model.RoleMenu{}

	// 删除原有的角色菜单关联
	role := &model.Role{}
	if err := role.DeleteRoleMenus(roleID); err != nil {
		return err
	}

	// 添加新的角色菜单关联
	if len(menuIDs) > 0 {
		var roleMenus []model.RoleMenu
		for _, menuID := range menuIDs {
			roleMenus = append(roleMenus, model.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			})
		}
		return roleMenu.BatchCreate(roleMenus)
	}

	return nil
}

// GetMenuTree 获取菜单树(用于前端菜单选择)
func (s *MenuService) GetMenuTree() ([]*model.Menu, error) {
	menu := &model.Menu{}
	return menu.GetAll()
}

// GetMenuPermsByUserID 获取用户菜单权限标识
func (s *MenuService) GetMenuPermsByUserID(userID uint) ([]string, error) {
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 获取用户角色
	casbinService := NewCasbinService()
	roles, err := casbinService.GetRolesForUser(userIDStr)
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
	menu := &model.Menu{}
	return menu.GetPermsByUserRoles(roleIDs)
}
