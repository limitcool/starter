package migration

import (
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/enum"
	"gorm.io/gorm"
)

// 这个文件用于添加API的迁移

func init() {
	// 注册API迁移
	RegisterMigration("016_init_apis", initAPIs, dropAPIs)
}

// initAPIs 初始化API
func initAPIs(tx *gorm.DB) error {
	// 检查是否已有API
	var count int64
	if err := tx.Model(&model.API{}).Count(&count).Error; err != nil {
		return err
	}

	// 已存在则不重复创建
	if count > 0 {
		return nil
	}

	// 创建基础API
	apis := []model.API{
		// 用户管理API
		{
			Path:        "/api/v1/admin/users",
			Method:      "GET",
			Name:        "获取用户列表",
			Description: "获取系统用户列表",
			Group:       "用户管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/users",
			Method:      "POST",
			Name:        "创建用户",
			Description: "创建系统用户",
			Group:       "用户管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/users/:id",
			Method:      "PUT",
			Name:        "更新用户",
			Description: "更新系统用户信息",
			Group:       "用户管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/users/:id",
			Method:      "DELETE",
			Name:        "删除用户",
			Description: "删除系统用户",
			Group:       "用户管理",
			Enabled:     true,
		},
		// 角色管理API
		{
			Path:        "/api/v1/admin/roles",
			Method:      "GET",
			Name:        "获取角色列表",
			Description: "获取角色列表",
			Group:       "角色管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/roles",
			Method:      "POST",
			Name:        "创建角色",
			Description: "创建角色",
			Group:       "角色管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/roles/:id",
			Method:      "PUT",
			Name:        "更新角色",
			Description: "更新角色信息",
			Group:       "角色管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/roles/:id",
			Method:      "DELETE",
			Name:        "删除角色",
			Description: "删除角色",
			Group:       "角色管理",
			Enabled:     true,
		},
		// 菜单管理API
		{
			Path:        "/api/v1/admin/menus",
			Method:      "GET",
			Name:        "获取菜单列表",
			Description: "获取菜单列表",
			Group:       "菜单管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/menus",
			Method:      "POST",
			Name:        "创建菜单",
			Description: "创建菜单",
			Group:       "菜单管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/menus/:id",
			Method:      "PUT",
			Name:        "更新菜单",
			Description: "更新菜单信息",
			Group:       "菜单管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/menus/:id",
			Method:      "DELETE",
			Name:        "删除菜单",
			Description: "删除菜单",
			Group:       "菜单管理",
			Enabled:     true,
		},
		// 权限管理API
		{
			Path:        "/api/v1/admin/permissions",
			Method:      "GET",
			Name:        "获取权限列表",
			Description: "获取权限列表",
			Group:       "权限管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/permissions",
			Method:      "POST",
			Name:        "创建权限",
			Description: "创建权限",
			Group:       "权限管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/permissions/:id",
			Method:      "PUT",
			Name:        "更新权限",
			Description: "更新权限信息",
			Group:       "权限管理",
			Enabled:     true,
		},
		{
			Path:        "/api/v1/admin/permissions/:id",
			Method:      "DELETE",
			Name:        "删除权限",
			Description: "删除权限",
			Group:       "权限管理",
			Enabled:     true,
		},
	}

	// 创建API
	for _, api := range apis {
		if err := tx.Create(&api).Error; err != nil {
			return err
		}

		// 创建对应的权限记录
		permission := model.Permission{
			Name:        api.Name,
			Code:        api.Path + ":" + api.Method,
			Type:        enum.PermissionTypeAPI,
			Description: api.Description,
			Enabled:     true,
			APIID:       api.ID,
		}

		if err := tx.Create(&permission).Error; err != nil {
			return err
		}
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

	var permMenu model.Menu
	if err := tx.Where("name = ? AND parent_id = ?", "权限管理", 1).First(&permMenu).Error; err != nil {
		return err
	}

	// 关联菜单和API
	menuAPIs := []struct {
		MenuID uint
		APIIDs []uint
	}{
		{
			MenuID: userMenu.ID,
			APIIDs: []uint{1, 2, 3, 4},
		},
		{
			MenuID: roleMenu.ID,
			APIIDs: []uint{5, 6, 7, 8},
		},
		{
			MenuID: menuMenu.ID,
			APIIDs: []uint{9, 10, 11, 12},
		},
		{
			MenuID: permMenu.ID,
			APIIDs: []uint{13, 14, 15, 16},
		},
	}

	for _, ma := range menuAPIs {
		for _, apiID := range ma.APIIDs {
			menuAPI := model.MenuAPI{
				MenuID: ma.MenuID,
				APIID:  apiID,
			}
			if err := tx.Create(&menuAPI).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// dropAPIs 删除API
func dropAPIs(tx *gorm.DB) error {
	// 删除菜单API关联
	if err := tx.Delete(&model.MenuAPI{}).Error; err != nil {
		return err
	}

	// 删除API对应的权限
	if err := tx.Where("type = ?", enum.PermissionTypeAPI).Delete(&model.Permission{}).Error; err != nil {
		return err
	}
	
	// 删除所有API
	return tx.Delete(&model.API{}).Error
}
