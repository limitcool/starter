package migration

import (
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/enum"
	"gorm.io/gorm"
)

// 这个文件用于添加菜单按钮的迁移

func init() {
	// 注册菜单按钮迁移
	RegisterMigration("015_init_menu_buttons", initMenuButtons, dropMenuButtons)
}

// initMenuButtons 初始化菜单按钮
func initMenuButtons(tx *gorm.DB) error {
	// 检查是否已有按钮
	var count int64
	if err := tx.Model(&model.MenuButton{}).Count(&count).Error; err != nil {
		return err
	}

	// 已存在则不重复创建
	if count > 0 {
		return nil
	}

	// 获取菜单
	var userMenu model.Menu
	if err := tx.Where("name = ? AND parent_id = ?", "用户管理", 1).First(&userMenu).Error; err != nil {
		return err
	}

	var roleMenu model.Menu
	if err := tx.Where("name = ? AND parent_id = ?", "角色管理", 1).First(&roleMenu).Error; err != nil {
		return err
	}

	var menuMenu model.Menu
	if err := tx.Where("name = ? AND parent_id = ?", "菜单管理", 1).First(&menuMenu).Error; err != nil {
		return err
	}

	// 创建按钮
	buttons := []model.MenuButton{
		// 用户管理按钮
		{
			MenuID:     userMenu.ID,
			Name:       "添加用户",
			Permission: "system:user:add",
			Icon:       "plus",
			OrderNum:   1,
			Enabled:    true,
		},
		{
			MenuID:     userMenu.ID,
			Name:       "编辑用户",
			Permission: "system:user:edit",
			Icon:       "edit",
			OrderNum:   2,
			Enabled:    true,
		},
		{
			MenuID:     userMenu.ID,
			Name:       "删除用户",
			Permission: "system:user:delete",
			Icon:       "delete",
			OrderNum:   3,
			Enabled:    true,
		},
		{
			MenuID:     userMenu.ID,
			Name:       "分配角色",
			Permission: "system:user:role",
			Icon:       "role",
			OrderNum:   4,
			Enabled:    true,
		},
		// 角色管理按钮
		{
			MenuID:     roleMenu.ID,
			Name:       "添加角色",
			Permission: "system:role:add",
			Icon:       "plus",
			OrderNum:   1,
			Enabled:    true,
		},
		{
			MenuID:     roleMenu.ID,
			Name:       "编辑角色",
			Permission: "system:role:edit",
			Icon:       "edit",
			OrderNum:   2,
			Enabled:    true,
		},
		{
			MenuID:     roleMenu.ID,
			Name:       "删除角色",
			Permission: "system:role:delete",
			Icon:       "delete",
			OrderNum:   3,
			Enabled:    true,
		},
		{
			MenuID:     roleMenu.ID,
			Name:       "分配权限",
			Permission: "system:role:permission",
			Icon:       "permission",
			OrderNum:   4,
			Enabled:    true,
		},
		// 菜单管理按钮
		{
			MenuID:     menuMenu.ID,
			Name:       "添加菜单",
			Permission: "system:menu:add",
			Icon:       "plus",
			OrderNum:   1,
			Enabled:    true,
		},
		{
			MenuID:     menuMenu.ID,
			Name:       "编辑菜单",
			Permission: "system:menu:edit",
			Icon:       "edit",
			OrderNum:   2,
			Enabled:    true,
		},
		{
			MenuID:     menuMenu.ID,
			Name:       "删除菜单",
			Permission: "system:menu:delete",
			Icon:       "delete",
			OrderNum:   3,
			Enabled:    true,
		},
	}

	// 创建按钮
	for _, button := range buttons {
		if err := tx.Create(&button).Error; err != nil {
			return err
		}

		// 创建对应的权限记录
		permission := model.Permission{
			Name:        button.Name,
			Code:        button.Permission,
			Type:        enum.PermissionTypeButton,
			Description: button.Name + "权限",
			Enabled:     true,
			ButtonID:    button.ID,
			MenuID:      button.MenuID,
		}

		if err := tx.Create(&permission).Error; err != nil {
			return err
		}
	}

	return nil
}

// dropMenuButtons 删除菜单按钮
func dropMenuButtons(tx *gorm.DB) error {
	// 删除按钮对应的权限
	if err := tx.Where("type = ?", enum.PermissionTypeButton).Delete(&model.Permission{}).Error; err != nil {
		return err
	}
	
	// 删除所有按钮
	return tx.Delete(&model.MenuButton{}).Error
}
