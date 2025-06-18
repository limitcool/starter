package migration

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/limitcool/starter/internal/model"
	"gorm.io/gorm"
)

// CreatePermissionTables20250618 创建权限相关表
func CreatePermissionTables20250618(db *gorm.DB) error {
	// 创建Casbin策略表
	if err := db.AutoMigrate(&gormadapter.CasbinRule{}); err != nil {
		return err
	}

	// 创建角色表
	if err := db.AutoMigrate(&model.Role{}); err != nil {
		return err
	}

	// 创建权限表
	if err := db.AutoMigrate(&model.Permission{}); err != nil {
		return err
	}

	// 创建菜单表
	if err := db.AutoMigrate(&model.Menu{}); err != nil {
		return err
	}

	// 注意：用户角色关联由Casbin管理，不需要user_roles表

	// 注意：在基于Casbin的设计中，不需要以下关联表：
	// - RolePermission: 角色权限关联由Casbin管理
	// - RoleMenu: 菜单通过permission_key与权限关联，不直接关联角色
	// - MenuPermission: 菜单通过permission_key字段关联权限

	return nil
}

// InitPermissionData20250618 初始化权限数据
func InitPermissionData20250618(db *gorm.DB) error {
	// 创建默认角色
	roles := []model.Role{
		{
			Name:        "超级管理员",
			Key:         "admin",
			Description: "系统管理员，拥有所有权限",
			Status:      1,
		},
		{
			Name:        "教练",
			Key:         "coach",
			Description: "管理课程和学员",
			Status:      1,
		},
		{
			Name:        "销售",
			Key:         "sales",
			Description: "管理会员和线索",
			Status:      1,
		},
		{
			Name:        "会员",
			Key:         "member",
			Description: "预约课程",
			Status:      1,
		},
	}

	for _, role := range roles {
		var existingRole model.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// 创建默认权限
	permissions := []model.Permission{
		// 系统管理权限分组
		{
			ParentID: 0,
			Name:     "系统管理",
			Key:      "sys",
			Type:     "MENU",
		},
		{
			ParentID: 1,
			Name:     "用户管理",
			Key:      "sys:user",
			Type:     "MENU",
		},
		{
			ParentID: 2,
			Name:     "查看用户列表",
			Key:      "user:list",
			Type:     "API",
		},
		{
			ParentID: 2,
			Name:     "创建用户",
			Key:      "user:create",
			Type:     "API",
		},
		// 会员管理权限
		{
			ParentID: 0,
			Name:     "会员管理",
			Key:      "member_manage",
			Type:     "MENU",
		},
		{
			ParentID: 5,
			Name:     "查看会员列表",
			Key:      "member:list",
			Type:     "API",
		},
		{
			ParentID: 5,
			Name:     "编辑会员信息",
			Key:      "member:edit",
			Type:     "API",
		},
		// 课程管理权限
		{
			ParentID: 0,
			Name:     "课程管理",
			Key:      "course_manage",
			Type:     "MENU",
		},
		{
			ParentID: 8,
			Name:     "查看课程列表",
			Key:      "course:list",
			Type:     "API",
		},
		{
			ParentID: 8,
			Name:     "创建课程",
			Key:      "course:create",
			Type:     "API",
		},
		// 小程序权限
		{
			ParentID: 0,
			Name:     "我的学员(小程序)",
			Key:      "mp_student:list",
			Type:     "API",
		},
	}

	for _, permission := range permissions {
		var existingPermission model.Permission
		if err := db.Where("name = ?", permission.Name).First(&existingPermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&permission).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// 创建默认菜单
	menus := []model.Menu{
		{
			ParentID:      0,
			Name:          "仪表盘",
			Path:          "/dashboard",
			Component:     "Dashboard",
			Icon:          "dashboard",
			SortOrder:     1,
			IsVisible:     true,
			PermissionKey: "", // 仪表盘不需要特殊权限
			Platform:      "admin",
		},
		{
			ParentID:      0,
			Name:          "系统管理",
			Path:          "/system",
			Component:     "Layout",
			Icon:          "system",
			SortOrder:     2,
			IsVisible:     true,
			PermissionKey: "sys", // 需要系统管理权限
			Platform:      "admin",
		},
		{
			ParentID:      2, // system菜单的ID
			Name:          "用户管理",
			Path:          "/system/user",
			Component:     "system/User",
			Icon:          "user",
			SortOrder:     1,
			IsVisible:     true,
			PermissionKey: "sys:user", // 需要用户管理权限
			Platform:      "admin",
		},
		{
			ParentID:      0,
			Name:          "会员管理",
			Path:          "/member",
			Component:     "member/Index",
			Icon:          "member",
			SortOrder:     3,
			IsVisible:     true,
			PermissionKey: "member_manage", // 需要会员管理权限
			Platform:      "admin",
		},
		{
			ParentID:      0,
			Name:          "课程管理",
			Path:          "/course",
			Component:     "course/Index",
			Icon:          "course",
			SortOrder:     4,
			IsVisible:     true,
			PermissionKey: "course_manage", // 需要课程管理权限
			Platform:      "admin",
		},
		{
			ParentID:      0,
			Name:          "我的学员",
			Path:          "/mp/students",
			Component:     "mp/Students",
			Icon:          "student",
			SortOrder:     1,
			IsVisible:     true,
			PermissionKey: "mp_student:list", // 需要学员列表权限
			Platform:      "coach_mp",
		},
	}

	for _, menu := range menus {
		var existingMenu model.Menu
		if err := db.Where("name = ?", menu.Name).First(&existingMenu).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&menu).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
