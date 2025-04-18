package model

import (
	"fmt"
	
	"gorm.io/gorm"
)

// GetMenuByID 根据ID获取菜单
func GetMenuByID(db *gorm.DB, id uint) (*Menu, error) {
	var menu Menu
	err := db.Where("id = ?", id).First(&menu).Error
	return &menu, err
}

// GetAllMenus 获取所有菜单
func GetAllMenus(db *gorm.DB) ([]*Menu, error) {
	var menus []*Menu
	err := db.Order("order_num").Find(&menus).Error
	return menus, err
}

// CreateMenu 创建菜单
func CreateMenu(db *gorm.DB, menu *Menu) error {
	return db.Create(menu).Error
}

// UpdateMenu 更新菜单
func UpdateMenu(db *gorm.DB, menu *Menu) error {
	return db.Model(&Menu{}).Where("id = ?", menu.ID).Updates(menu).Error
}

// DeleteMenu 删除菜单
func DeleteMenu(db *gorm.DB, id uint) error {
	// 检查是否有子菜单
	var count int64
	if err := db.Model(&Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该菜单下有子菜单，无法删除")
	}

	// 删除菜单
	return db.Delete(&Menu{}, id).Error
}

// GetMenusByRoleID 获取角色菜单
func GetMenusByRoleID(db *gorm.DB, roleID uint) ([]*Menu, error) {
	// 查询角色关联的菜单ID
	var roleMenus []RoleMenu
	err := db.Where("role_id = ?", roleID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	if len(menuIDs) == 0 {
		return []*Menu{}, nil
	}

	// 查询菜单
	var menus []*Menu
	err = db.Where("id IN ?", menuIDs).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

// GetMenusByUserID 获取用户菜单
func GetMenusByUserID(db *gorm.DB, userID uint) ([]*Menu, error) {
	// 1. 获取用户角色
	var userRoles []UserRole
	err := db.Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	// 提取角色ID
	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	if len(roleIDs) == 0 {
		return []*Menu{}, nil
	}

	// 2. 获取角色关联的菜单ID
	var roleMenus []RoleMenu
	err = db.Where("role_id IN ?", roleIDs).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	if len(menuIDs) == 0 {
		return []*Menu{}, nil
	}

	// 3. 查询菜单信息
	var menus []*Menu
	err = db.Where("id IN ? AND status = ? AND type IN ?", menuIDs, 1, []int8{0, 1}).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

// GetMenuPermsByUserID 获取用户菜单权限
func GetMenuPermsByUserID(db *gorm.DB, userID uint) ([]string, error) {
	// 获取用户角色
	var userRoles []UserRole
	err := db.Where("user_id = ?", userID).Find(&userRoles).Error
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

	return GetMenuPermsByRoleIDs(db, roleIDs)
}

// GetMenuPermsByRoleIDs 获取角色菜单权限
func GetMenuPermsByRoleIDs(db *gorm.DB, roleIDs []uint) ([]string, error) {
	// 查询角色菜单
	var roleMenus []RoleMenu
	err := db.Where("role_id IN ?", roleIDs).Find(&roleMenus).Error
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
	err = db.Model(&Menu{}).
		Where("id IN ? AND status = ? AND perms != ''", menuIDs, 1).
		Pluck("perms", &perms).Error
	if err != nil {
		return nil, err
	}

	return perms, nil
}
