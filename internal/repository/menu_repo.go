package repository

import (
	"sort"

	"github.com/limitcool/starter/internal/model"
	"gorm.io/gorm"
)

// MenuRepo 菜单仓库
type MenuRepo struct {
	DB *gorm.DB
}

// NewMenuRepo 创建菜单仓库
func NewMenuRepo(db *gorm.DB) *MenuRepo {
	return &MenuRepo{DB: db}
}

// GetByID 根据ID获取菜单
func (r *MenuRepo) GetByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.DB.Where("id = ?", id).First(&menu).Error
	return &menu, err
}

// GetByIDWithRelations 根据ID获取菜单及其关联数据
func (r *MenuRepo) GetByIDWithRelations(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.DB.Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}

	// 加载按钮
	r.DB.Where("menu_id = ?", menu.ID).Find(&menu.Buttons)

	// 加载API
	r.DB.Model(&menu).Association("APIs").Find(&menu.APIs)

	return &menu, nil
}

// GetAll 获取所有菜单
func (r *MenuRepo) GetAll() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.DB.Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 加载按钮
	for _, menu := range menus {
		r.DB.Where("menu_id = ?", menu.ID).Find(&menu.Buttons)
	}

	return menus, err
}

// Create 创建菜单
func (r *MenuRepo) Create(menu *model.Menu) error {
	return r.DB.Create(menu).Error
}

// CreateButton 创建菜单按钮
func (r *MenuRepo) CreateButton(button *model.MenuButton) error {
	return r.DB.Create(button).Error
}

// UpdateButton 更新菜单按钮
func (r *MenuRepo) UpdateButton(button *model.MenuButton) error {
	return r.DB.Save(button).Error
}

// DeleteButton 删除菜单按钮
func (r *MenuRepo) DeleteButton(id uint) error {
	return r.DB.Delete(&model.MenuButton{}, id).Error
}

// GetButtonByID 根据ID获取菜单按钮
func (r *MenuRepo) GetButtonByID(id uint) (*model.MenuButton, error) {
	var button model.MenuButton
	err := r.DB.First(&button, id).Error
	return &button, err
}

// Update 更新菜单
func (r *MenuRepo) Update(menu *model.Menu) error {
	return r.DB.Model(&model.Menu{}).Where("id = ?", menu.ID).Updates(menu).Error
}

// Delete 删除菜单
func (r *MenuRepo) Delete(id uint) error {
	// 检查是否有子菜单
	var count int64
	if err := r.DB.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrForeignKeyViolated
	}

	// 删除菜单
	return r.DB.Delete(&model.Menu{}, id).Error
}

// GetByRoleID 获取角色菜单
func (r *MenuRepo) GetByRoleID(roleID uint) ([]*model.Menu, error) {
	// 查询角色关联的菜单ID
	var roleMenus []model.RoleMenu
	err := r.DB.Where("role_id = ?", roleID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	if len(menuIDs) == 0 {
		return []*model.Menu{}, nil
	}

	// 查询菜单
	var menus []*model.Menu
	err = r.DB.Where("id IN ?", menuIDs).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

// GetByUserID 获取用户菜单
func (r *MenuRepo) GetByUserID(userID uint) ([]*model.Menu, error) {
	// 1. 获取用户角色
	var userRoles []model.UserRole
	err := r.DB.Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	// 提取角色ID
	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	if len(roleIDs) == 0 {
		return []*model.Menu{}, nil
	}

	// 2. 获取角色关联的菜单ID
	var roleMenus []model.RoleMenu
	err = r.DB.Where("role_id IN ?", roleIDs).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	if len(menuIDs) == 0 {
		return []*model.Menu{}, nil
	}

	// 3. 查询菜单信息
	var menus []*model.Menu
	err = r.DB.Where("id IN ? AND status = ? AND type IN ?", menuIDs, 1, []int8{0, 1}).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, err
}

// GetPermsByUserID 获取用户菜单权限
func (r *MenuRepo) GetPermsByUserID(userID uint) ([]string, error) {
	// 获取用户角色
	var userRoles []model.UserRole
	err := r.DB.Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	// 提取角色ID
	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	return r.GetPermsByRoleIDs(roleIDs)
}

// GetPermsByRoleIDs 获取角色菜单权限
func (r *MenuRepo) GetPermsByRoleIDs(roleIDs []uint) ([]string, error) {
	// 查询角色菜单
	var roleMenus []model.RoleMenu
	err := r.DB.Where("role_id IN ?", roleIDs).Find(&roleMenus).Error
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
	err = r.DB.Model(&model.Menu{}).
		Where("id IN ? AND status = ? AND perms != ''", menuIDs, 1).
		Pluck("perms", &perms).Error
	if err != nil {
		return nil, err
	}

	return perms, nil
}

// AssociateAPI 关联菜单和API
func (r *MenuRepo) AssociateAPI(menuID uint, apiIDs []uint) error {
	// 开始事务
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 获取菜单
	var menu model.Menu
	if err := tx.First(&menu, menuID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除原有关联
	if err := tx.Where("menu_id = ?", menuID).Delete(&model.MenuAPI{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新关联
	if len(apiIDs) > 0 {
		var menuAPIs []model.MenuAPI
		for _, apiID := range apiIDs {
			menuAPIs = append(menuAPIs, model.MenuAPI{
				MenuID: menuID,
				APIID:  apiID,
			})
		}
		if err := tx.Create(&menuAPIs).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// AssignMenuToRole 为角色分配菜单
func (r *MenuRepo) AssignMenuToRole(roleID uint, menuIDs []uint) error {
	// 开始事务
	tx := r.DB.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 删除原有的角色菜单关联
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		tx.Rollback()
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
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// BuildMenuTree 构建菜单树
func (r *MenuRepo) BuildMenuTree(menus []*model.Menu) []*model.MenuTree {
	// 创建一个映射，用于快速查找菜单
	menuMap := make(map[uint]*model.MenuTree)

	// 将所有菜单转换为树节点
	for _, menu := range menus {
		menuMap[menu.ID] = &model.MenuTree{
			Menu:     menu,
			Children: []*model.MenuTree{},
		}
	}

	// 构建树结构
	var rootMenus []*model.MenuTree
	for _, menu := range menus {
		if menu.ParentID == 0 {
			// 根菜单
			rootMenus = append(rootMenus, menuMap[menu.ID])
		} else {
			// 子菜单
			if parent, ok := menuMap[menu.ParentID]; ok {
				parent.Children = append(parent.Children, menuMap[menu.ID])
			}
		}
	}

	return rootMenus
}

// BuildMenuTreeOld 构建菜单树(旧版本)
func (r *MenuRepo) BuildMenuTreeOld(menus []*model.Menu) []*model.Menu {
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

// AddMenuAPI 添加菜单API关联
func (r *MenuRepo) AddMenuAPI(menuID uint, apiID uint) error {
	menuAPI := model.MenuAPI{
		MenuID: menuID,
		APIID:  apiID,
	}
	return r.DB.Create(&menuAPI).Error
}

// ClearMenuAPIs 清除菜单API关联
func (r *MenuRepo) ClearMenuAPIs(menuID uint) error {
	return r.DB.Where("menu_id = ?", menuID).Delete(&model.MenuAPI{}).Error
}

// GetMenuIDsByAPIID 获取API关联的所有菜单ID
func (r *MenuRepo) GetMenuIDsByAPIID(apiID uint) ([]uint, error) {
	var menuAPIs []model.MenuAPI
	err := r.DB.Where("api_id = ?", apiID).Find(&menuAPIs).Error
	if err != nil {
		return nil, err
	}

	var menuIDs []uint
	for _, ma := range menuAPIs {
		menuIDs = append(menuIDs, ma.MenuID)
	}

	return menuIDs, nil
}

// GetRolesByMenuID 获取拥有该菜单的所有角色
func (r *MenuRepo) GetRolesByMenuID(menuID uint) ([]*model.Role, error) {
	var roleMenus []model.RoleMenu
	err := r.DB.Where("menu_id = ?", menuID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	var roleIDs []uint
	for _, rm := range roleMenus {
		roleIDs = append(roleIDs, rm.RoleID)
	}

	if len(roleIDs) == 0 {
		return []*model.Role{}, nil
	}

	var roles []*model.Role
	err = r.DB.Where("id IN ?", roleIDs).Find(&roles).Error
	return roles, err
}
