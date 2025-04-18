package repository

import (
	"github.com/limitcool/starter/internal/model"
	"gorm.io/gorm"
)

// MenuRepo 菜单仓库
type MenuRepo struct {
	DB *gorm.DB
}

// NewMenuRepo 创建菜单仓库
func NewMenuRepo(db *gorm.DB) MenuRepository {
	return &MenuRepo{DB: db}
}

// GetByID 根据ID获取菜单
func (r *MenuRepo) GetByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.DB.Where("id = ?", id).First(&menu).Error
	return &menu, err
}

// GetAll 获取所有菜单
func (r *MenuRepo) GetAll() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.DB.Order("order_num").Find(&menus).Error
	return menus, err
}

// Create 创建菜单
func (r *MenuRepo) Create(menu *model.Menu) error {
	return r.DB.Create(menu).Error
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
