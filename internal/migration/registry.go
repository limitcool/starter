package migration

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/enum"
	"gorm.io/gorm"
)

// 注册所有迁移
// 这个文件用于注册所有迁移
// 在实际项目中，你可以将迁移分散到多个文件中，每个文件对应一个模块

func init() {
	// 基础表迁移
	RegisterMigration("001_create_file_table", createFileTable, dropFileTable)

	// 用户表迁移
	RegisterMigration("002_create_sys_users_table", createSysUsersTable, dropSysUsersTable)
	RegisterMigration("003_create_users_table", createUsersTable, dropUsersTable)

	// 角色表迁移
	RegisterMigration("004_create_roles_tables", createRolesTables, dropRolesTables)
	RegisterMigration("005_init_roles", initRoles, dropRoles)

	// API和菜单表迁移
	RegisterMigration("006_create_api_table", createAPITable, dropAPITable)
	RegisterMigration("007_create_menus_table", createMenusTable, dropMenusTable)
	RegisterMigration("008_create_menu_button_table", createMenuButtonTable, dropMenuButtonTable)
	RegisterMigration("009_create_menu_api_table", createMenuAPITable, dropMenuAPITable)

	// 权限表迁移
	RegisterMigration("010_create_permissions_tables", createPermissionsTables, dropPermissionsTables)
	RegisterMigration("011_create_casbin_rule_table", createCasbinRuleTable, dropCasbinRuleTable)

	// 操作日志表迁移
	RegisterMigration("012_create_operation_logs_table", createOperationLogsTable, dropOperationLogsTable)

	// 初始数据迁移
	RegisterMigration("013_init_basic_menus", initBasicMenus, dropBasicMenus)
	RegisterMigration("014_init_admin_user",
		// 包装函数，使其符合 func(*gorm.DB) error 签名
		func(tx *gorm.DB) error {
			return initAdminUser(tx, nil)
		},
		dropAdminUser)
}

// 文件表迁移
func createFileTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.File{})
}

func dropFileTable(tx *gorm.DB) error {
	return tx.Migrator().DropTable("file")
}

// 用户表迁移
func createSysUsersTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.SysUser{})
}

func dropSysUsersTable(tx *gorm.DB) error {
	return tx.Migrator().DropTable("sys_user")
}

func createUsersTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.User{})
}

func dropUsersTable(tx *gorm.DB) error {
	return tx.Migrator().DropTable("user")
}

// 角色表迁移
func createRolesTables(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&model.Role{}); err != nil {
		return err
	}
	if err := tx.AutoMigrate(&model.UserRole{}); err != nil {
		return err
	}
	if err := tx.AutoMigrate(&model.RoleMenu{}); err != nil {
		return err
	}
	return nil
}

func dropRolesTables(tx *gorm.DB) error {
	if err := tx.Migrator().DropTable("sys_user_role"); err != nil {
		return err
	}
	if err := tx.Migrator().DropTable("sys_role_menu"); err != nil {
		return err
	}
	return tx.Migrator().DropTable("sys_role")
}

func initRoles(tx *gorm.DB) error {
	// 检查是否已有admin角色
	var count int64
	if err := tx.Model(&model.Role{}).Where("code = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	// 已存在则不重复创建
	if count > 0 {
		return nil
	}

	// 创建超级管理员角色
	adminRole := model.Role{
		Name:        "超级管理员",
		Code:        "admin",
		Enabled:     true,
		Sort:        0,
		Description: "系统超级管理员",
	}

	return tx.Create(&adminRole).Error
}

func dropRoles(tx *gorm.DB) error {
	return tx.Where("code = ?", "admin").Delete(&model.Role{}).Error
}

// 权限表迁移
func createPermissionsTables(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&model.Permission{}); err != nil {
		return err
	}
	if err := tx.AutoMigrate(&model.RolePermission{}); err != nil {
		return err
	}
	return nil
}

func dropPermissionsTables(tx *gorm.DB) error {
	if err := tx.Migrator().DropTable("sys_role_permission"); err != nil {
		return err
	}
	return tx.Migrator().DropTable("sys_permission")
}

func createCasbinRuleTable(tx *gorm.DB) error {
	// 检查表是否存在
	if tx.Migrator().HasTable("casbin_rule") {
		return nil
	}

	// 获取数据库类型
	var dbType string
	if tx.Dialector.Name() == "mysql" || tx.Dialector.Name() == "sqlite" {
		// MySQL或SQLite语法
		dbType = "mysql"
	} else {
		// PostgreSQL语法
		dbType = "postgres"
	}

	// 根据数据库类型选择不同的SQL语句
	var sql string
	if dbType == "mysql" {
		// MySQL语法
		sql = `CREATE TABLE casbin_rule (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			ptype VARCHAR(100) DEFAULT NULL,
			v0 VARCHAR(100) DEFAULT NULL,
			v1 VARCHAR(100) DEFAULT NULL,
			v2 VARCHAR(100) DEFAULT NULL,
			v3 VARCHAR(100) DEFAULT NULL,
			v4 VARCHAR(100) DEFAULT NULL,
			v5 VARCHAR(100) DEFAULT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY idx_casbin_rule (ptype,v0,v1,v2,v3,v4,v5)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	} else {
		// PostgreSQL语法
		sql = `CREATE TABLE casbin_rule (
			id BIGSERIAL PRIMARY KEY,
			ptype VARCHAR(100) DEFAULT NULL,
			v0 VARCHAR(100) DEFAULT NULL,
			v1 VARCHAR(100) DEFAULT NULL,
			v2 VARCHAR(100) DEFAULT NULL,
			v3 VARCHAR(100) DEFAULT NULL,
			v4 VARCHAR(100) DEFAULT NULL,
			v5 VARCHAR(100) DEFAULT NULL,
			UNIQUE (ptype,v0,v1,v2,v3,v4,v5)
		);`
	}

	return tx.Exec(sql).Error
}

func dropCasbinRuleTable(tx *gorm.DB) error {
	return tx.Migrator().DropTable("casbin_rule")
}

// API表迁移
func createAPITable(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.API{})
}

func dropAPITable(tx *gorm.DB) error {
	return tx.Migrator().DropTable("sys_api")
}

// 菜单表迁移
func createMenusTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.Menu{})
}

func dropMenusTable(tx *gorm.DB) error {
	return tx.Migrator().DropTable("sys_menu")
}

// 菜单按钮表迁移
func createMenuButtonTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.MenuButton{})
}

func dropMenuButtonTable(tx *gorm.DB) error {
	return tx.Migrator().DropTable("sys_menu_button")
}

// 菜单API关联表迁移
func createMenuAPITable(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.MenuAPI{})
}

func dropMenuAPITable(tx *gorm.DB) error {
	return tx.Migrator().DropTable("sys_menu_api")
}

// 操作日志表迁移
func createOperationLogsTable(tx *gorm.DB) error {
	return tx.AutoMigrate(&model.OperationLog{})
}

func dropOperationLogsTable(tx *gorm.DB) error {
	return tx.Migrator().DropTable("sys_operation_log")
}

// 初始数据迁移
func initBasicMenus(tx *gorm.DB) error {
	// 检查是否已有菜单
	var count int64
	if err := tx.Model(&model.Menu{}).Count(&count).Error; err != nil {
		return err
	}

	// 已存在则不重复创建
	if count > 0 {
		return nil
	}

	// 创建基础菜单
	menus := []model.Menu{
		{
			Name:      "系统管理",
			Path:      "/system",
			Component: "Layout",
			Redirect:  "/system/user",
			Icon:      "setting",
			Title:     "系统管理",
			OrderNum:  100,
			Type:      enum.MenuTypeDirectory, // 目录类型
			Enabled:   true,
			ParentID:  0,
		},
		{
			Name:      "用户管理",
			Path:      "user",
			Component: "/system/user/index",
			Icon:      "user",
			Title:     "用户管理",
			OrderNum:  101,
			Perms:     "system:user:list",
			Type:      enum.MenuTypeMenu, // 菜单类型
			Enabled:   true,
			ParentID:  1, // 系统管理的ID
		},
		{
			Name:      "角色管理",
			Path:      "role",
			Component: "/system/role/index",
			Icon:      "role",
			Title:     "角色管理",
			OrderNum:  102,
			Perms:     "system:role:list",
			Type:      enum.MenuTypeMenu, // 菜单类型
			Enabled:   true,
			ParentID:  1, // 系统管理的ID
		},
		{
			Name:      "菜单管理",
			Path:      "menu",
			Component: "/system/menu/index",
			Icon:      "menu",
			Title:     "菜单管理",
			OrderNum:  103,
			Perms:     "system:menu:list",
			Type:      enum.MenuTypeMenu, // 菜单类型
			Enabled:   true,
			ParentID:  1, // 系统管理的ID
		},
		{
			Name:      "权限管理",
			Path:      "permission",
			Component: "/system/permission/index",
			Icon:      "permission",
			Title:     "权限管理",
			OrderNum:  104,
			Perms:     "system:permission:list",
			Type:      enum.MenuTypeMenu, // 菜单类型
			Enabled:   true,
			ParentID:  1, // 系统管理的ID
		},
		{
			Name:      "API管理",
			Path:      "api",
			Component: "/system/api/index",
			Icon:      "api",
			Title:     "API管理",
			OrderNum:  105,
			Perms:     "system:api:list",
			Type:      enum.MenuTypeMenu, // 菜单类型
			Enabled:   true,
			ParentID:  1, // 系统管理的ID
		},
		{
			Name:      "日志管理",
			Path:      "log",
			Component: "/system/log/index",
			Icon:      "log",
			Title:     "日志管理",
			OrderNum:  106,
			Perms:     "system:log:list",
			Type:      enum.MenuTypeMenu, // 菜单类型
			Enabled:   true,
			ParentID:  1, // 系统管理的ID
		},
	}

	// 创建菜单
	for i := range menus {
		if err := tx.Create(&menus[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

func dropBasicMenus(tx *gorm.DB) error {
	return tx.Where("id IN (?)", []int{1, 2, 3, 4, 5, 6, 7}).Delete(&model.Menu{}).Error
}

func initAdminUser(tx *gorm.DB, config ...*configs.Config) error {
	// 获取配置
	var cfg *configs.Config
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	} else {
		cfg = &configs.Config{
			Admin: configs.Admin{
				Username: "admin",
				Password: "123456",
				Nickname: "超级管理员",
			},
		}
	}

	// 如果配置文件中没有设置管理员信息，使用默认值
	username := cfg.Admin.Username
	password := cfg.Admin.Password
	nickname := cfg.Admin.Nickname

	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "123456"
	}
	if nickname == "" {
		nickname = "超级管理员"
	}

	// 检查是否已有管理员用户
	var count int64
	if err := tx.Model(&model.AdminUser{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return err
	}

	// 已存在则不重复创建
	if count > 0 {
		return nil
	}

	// 获取管理员角色
	var adminRole model.Role
	if err := tx.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// 创建管理员用户
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}

	// 使用雪花ID
	adminUser := model.AdminUser{
		Username: username,
		Password: hashedPassword,
		Nickname: nickname,
		Email:    "admin@example.com",
		Enabled:  true,
	}

	if err := tx.Create(&adminUser).Error; err != nil {
		return err
	}

	// 关联角色
	userRole := struct {
		AdminUserID int64 `gorm:"column:admin_user_id"`
		RoleID      int64 `gorm:"column:role_id"`
	}{
		AdminUserID: adminUser.ID,
		RoleID:      int64(adminRole.ID),
	}

	return tx.Table("admin_user_role").Create(&userRole).Error
}

func dropAdminUser(tx *gorm.DB) error {
	// 获取管理员用户
	var adminUser model.AdminUser
	if err := tx.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		return nil // 如果不存在，直接返回
	}

	// 删除用户角色关联
	if err := tx.Table("admin_user_role").Where("admin_user_id = ?", adminUser.ID).Delete(nil).Error; err != nil {
		return err
	}

	// 删除用户
	return tx.Delete(&adminUser).Error
}
