package services

import (
	"errors"
	"sort"
	"strconv"

	"github.com/limitcool/starter/internal/model"
	"gorm.io/gorm"
)

// MenuService 菜单服务
type MenuService struct {
	db *gorm.DB
}

// NewMenuService 创建菜单服务
func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{
		db: db,
	}
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(menu *model.Menu) error {
	return s.db.Create(menu).Error
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(menu *model.Menu) error {
	return s.db.Model(&model.Menu{}).Where("id = ?", menu.ID).Updates(menu).Error
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	// 检查是否有子菜单
	var count int64
	if err := s.db.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该菜单下有子菜单，不能删除")
	}

	// 删除菜单关联的角色
	if err := s.db.Where("menu_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
		return err
	}

	// 删除菜单
	return s.db.Delete(&model.Menu{}, id).Error
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := s.db.Where("id = ?", id).First(&menu).Error
	return &menu, err
}

// GetAllMenus 获取所有菜单
func (s *MenuService) GetAllMenus() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := s.db.Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return s.buildMenuTree(menus), nil
}

// GetMenusByRoleID 获取角色菜单
func (s *MenuService) GetMenusByRoleID(roleID uint) ([]*model.Menu, error) {
	// 查询角色关联的菜单ID
	var roleMenus []model.RoleMenu
	err := s.db.Where("role_id = ?", roleID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	// 如果没有菜单，返回空数组
	if len(menuIDs) == 0 {
		return []*model.Menu{}, nil
	}

	// 查询菜单
	var menus []*model.Menu
	err = s.db.Where("id IN ?", menuIDs).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return s.buildMenuTree(menus), nil
}

// GetUserMenus 获取用户菜单
func (s *MenuService) GetUserMenus(userID uint) ([]*model.Menu, error) {
	// 1. 获取用户角色
	var userRoles []model.UserRole
	err := s.db.Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	// 如果用户没有角色，返回空数组
	if len(userRoles) == 0 {
		return []*model.Menu{}, nil
	}

	// 提取角色ID
	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	// 2. 获取角色关联的菜单ID
	var roleMenus []model.RoleMenu
	err = s.db.Where("role_id IN ?", roleIDs).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取菜单ID并去重
	menuIDMap := make(map[uint]bool)
	for _, rm := range roleMenus {
		menuIDMap[rm.MenuID] = true
	}

	if len(menuIDMap) == 0 {
		return []*model.Menu{}, nil
	}

	var menuIDs []uint
	for id := range menuIDMap {
		menuIDs = append(menuIDs, id)
	}

	// 3. 查询菜单信息
	var menus []*model.Menu
	err = s.db.Where("id IN ? AND status = ? AND type IN ?", menuIDs, 1, []int8{0, 1}).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 4. 构建菜单树
	return s.buildMenuTree(menus), nil
}

// 构建菜单树
func (s *MenuService) buildMenuTree(menus []*model.Menu) []*model.Menu {
	// 创建一个map用于快速查找
	menuMap := make(map[uint]*model.Menu)
	for _, m := range menus {
		menuMap[m.ID] = m
	}

	var rootMenus []*model.Menu
	for _, m := range menus {
		if m.ParentID == 0 {
			// 顶级菜单
			rootMenus = append(rootMenus, m)
		} else {
			// 子菜单
			if parent, ok := menuMap[m.ParentID]; ok {
				if parent.Children == nil {
					parent.Children = []*model.Menu{}
				}
				parent.Children = append(parent.Children, m)
			}
		}
	}

	// 菜单排序
	for _, m := range menuMap {
		if len(m.Children) > 0 {
			sort.Slice(m.Children, func(i, j int) bool {
				return m.Children[i].OrderNum < m.Children[j].OrderNum
			})
		}
	}

	// 对根菜单排序
	sort.Slice(rootMenus, func(i, j int) bool {
		return rootMenus[i].OrderNum < rootMenus[j].OrderNum
	})

	return rootMenus
}

// 为角色分配菜单
func (s *MenuService) AssignMenuToRole(roleID uint, menuIDs []uint) error {
	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除原有的角色菜单关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
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
			if err := tx.Create(&roleMenus).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetMenuTree 获取菜单树(用于前端菜单选择)
func (s *MenuService) GetMenuTree() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := s.db.Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return s.buildMenuTree(menus), nil
}

// GetMenusByUserID 获取用户菜单权限标识
func (s *MenuService) GetMenuPermsByUserID(userID uint) ([]string, error) {
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 获取用户角色
	casbinService := NewCasbinService(s.db)
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

	// 查询角色菜单
	var roleMenus []model.RoleMenu
	err = s.db.Where("role_id IN ?", roleIDs).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	if len(menuIDs) == 0 {
		return []string{}, nil
	}

	// 查询菜单权限标识
	var perms []string
	err = s.db.Model(&model.Menu{}).
		Where("id IN ? AND status = ? AND perms != ''", menuIDs, 1).
		Pluck("perms", &perms).Error
	if err != nil {
		return nil, err
	}

	return perms, nil
}
