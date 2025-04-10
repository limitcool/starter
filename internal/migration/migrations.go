package migration

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"gorm.io/gorm"
)

// 此文件用于定义具体的迁移实例
// 您可以在这里添加新的迁移，按照时间顺序排列

// RegisterCoreUserMigrations 注册用户相关迁移
func RegisterCoreUserMigrations(migrator *Migrator) {
	// 系统用户表迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080001",
		Name:    "create_sys_users_table",
		Up: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.SysUser{})
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("sys_user")
		},
	})

	// 普通用户表迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080007",
		Name:    "create_users_table",
		Up: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.User{})
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("user")
		},
	})
}

// RegisterRoleMigrations 注册角色相关迁移
func RegisterRoleMigrations(migrator *Migrator) {
	// 角色表迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080002",
		Name:    "create_roles_tables",
		Up: func(tx *gorm.DB) error {
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
		},
		Down: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable("sys_user_role"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("sys_role_menu"); err != nil {
				return err
			}
			return tx.Migrator().DropTable("sys_role")
		},
	})

	// 创建基础角色
	migrator.Register(&MigrationEntry{
		Version: "202504080003",
		Name:    "init_roles",
		Up: func(tx *gorm.DB) error {
			// 检查是否已有admin角色
			var count int64
			if err := tx.Model(&model.Role{}).Where("code = ?", "admin").Count(&count).Error; err != nil {
				return err
			}

			// 已存在则不重复创建
			if count > 0 {
				log.Info("管理员角色已存在，跳过创建")
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

			if err := tx.Create(&adminRole).Error; err != nil {
				return err
			}

			log.Info("管理员角色创建成功")
			return nil
		},
		Down: func(tx *gorm.DB) error {
			return tx.Where("code = ?", "admin").Delete(&model.Role{}).Error
		},
	})
}

// RegisterPermissionMigrations 注册权限相关迁移
func RegisterPermissionMigrations(migrator *Migrator) {
	// 权限表迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080004",
		Name:    "create_permissions_tables",
		Up: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&model.Permission{}); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&model.RolePermission{}); err != nil {
				return err
			}
			return nil
		},
		Down: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable("sys_role_permission"); err != nil {
				return err
			}
			return tx.Migrator().DropTable("sys_permission")
		},
	})
}

// RegisterMenuMigrations 注册菜单相关迁移
func RegisterMenuMigrations(migrator *Migrator) {
	// 菜单表迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080005",
		Name:    "create_menus_table",
		Up: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Menu{})
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("sys_menu")
		},
	})
}

// RegisterOperationLogMigrations 注册操作日志相关迁移
func RegisterOperationLogMigrations(migrator *Migrator) {
	// 操作日志表迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080006",
		Name:    "create_operation_logs_table",
		Up: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.OperationLog{})
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("sys_operation_log")
		},
	})
}

// RegisterInitialDataMigrations 注册初始数据迁移
func RegisterInitialDataMigrations(migrator *Migrator) {
	// 初始化菜单数据
	migrator.Register(&MigrationEntry{
		Version: "202504080008",
		Name:    "init_basic_menus",
		Up: func(tx *gorm.DB) error {
			// 检查是否已有菜单
			var count int64
			if err := tx.Model(&model.Menu{}).Count(&count).Error; err != nil {
				return err
			}

			// 已存在则不重复创建
			if count > 0 {
				log.Info("基础菜单已存在，跳过创建")
				return nil
			}

			// 获取管理员角色
			var adminRole model.Role
			if err := tx.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
				log.Error("获取管理员角色失败，跳过菜单创建", "error", err)
				return nil
			}

			// 创建系统管理菜单
			sysManage := model.Menu{
				Name:      "系统管理",
				ParentID:  0,
				Path:      "/system",
				Component: "Layout",
				Type:      0, // 目录
				Icon:      "system",
				OrderNum:  1,
				IsFrame:   false,
				IsHidden:  false,
				Enabled:   true,
				Perms:     "",
			}

			if err := tx.Create(&sysManage).Error; err != nil {
				return err
			}

			// 创建用户管理菜单
			userManage := model.Menu{
				Name:      "用户管理",
				ParentID:  sysManage.ID,
				Path:      "user",
				Component: "system/user/index",
				Type:      1, // 菜单
				Icon:      "user",
				OrderNum:  1,
				IsFrame:   false,
				IsHidden:  false,
				Enabled:   true,
				Perms:     "system:user:list",
			}

			if err := tx.Create(&userManage).Error; err != nil {
				return err
			}

			// 创建角色管理菜单
			roleManage := model.Menu{
				Name:      "角色管理",
				ParentID:  sysManage.ID,
				Path:      "role",
				Component: "system/role/index",
				Type:      1, // 菜单
				Icon:      "role",
				OrderNum:  2,
				IsFrame:   false,
				IsHidden:  false,
				Enabled:   true,
				Perms:     "system:role:list",
			}

			if err := tx.Create(&roleManage).Error; err != nil {
				return err
			}

			// 创建菜单管理菜单
			menuManage := model.Menu{
				Name:      "菜单管理",
				ParentID:  sysManage.ID,
				Path:      "menu",
				Component: "system/menu/index",
				Type:      1, // 菜单
				Icon:      "menu",
				OrderNum:  3,
				IsFrame:   false,
				IsHidden:  false,
				Enabled:   true,
				Perms:     "system:menu:list",
			}

			if err := tx.Create(&menuManage).Error; err != nil {
				return err
			}

			// 为超级管理员分配所有菜单
			allMenus := []model.Menu{sysManage, userManage, roleManage, menuManage}
			for _, menu := range allMenus {
				roleMenu := model.RoleMenu{
					RoleID: adminRole.ID,
					MenuID: menu.ID,
				}

				if err := tx.Create(&roleMenu).Error; err != nil {
					return err
				}
			}

			log.Info("基础菜单创建成功")
			return nil
		},
		Down: func(tx *gorm.DB) error {
			// 删除菜单
			return tx.Exec("DELETE FROM sys_menu").Error
		},
	})

	// 初始化管理员账号
	migrator.Register(&MigrationEntry{
		Version: "202504080009",
		Name:    "init_admin_user",
		Up: func(tx *gorm.DB) error {
			// 获取配置
			cfg := migrator.config
			if cfg == nil {
				log.Warn("配置未初始化，使用默认管理员账号")
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
			if err := tx.Model(&model.SysUser{}).Where("username = ?", username).Count(&count).Error; err != nil {
				return err
			}

			// 已存在则跳过
			if count > 0 {
				log.Info("管理员用户已存在，跳过创建")
				return nil
			}

			// 创建超级管理员用户
			hashedPassword, err := crypto.HashPassword(password)
			if err != nil {
				return fmt.Errorf("密码加密失败: %w", err)
			}

			// 使用GORM模型创建管理员用户，确保触发雪花ID生成
			sysUser := &model.SysUser{
				Username:     username,
				Password:     hashedPassword,
				Nickname:     nickname,
				Enabled:      true,
				Remark:       "系统初始化创建",
				AvatarFileID: 0,
				Email:        "",
				Mobile:       "",
				LastLogin:    nil, // 使用nil
				LastIP:       "",
			}

			log.Info("准备创建管理员用户",
				"username", sysUser.Username,
				"nickname", sysUser.Nickname)

			if err := tx.Create(sysUser).Error; err != nil {
				return fmt.Errorf("创建管理员账号失败: %w", err)
			}

			// 获取管理员角色ID
			var role model.Role
			if err := tx.Where("code = ?", "admin").First(&role).Error; err != nil {
				log.Warn("获取管理员角色失败，跳过角色分配", "error", err)
				return nil
			}

			// 创建用户角色关联
			userRole := &model.UserRole{
				UserID: sysUser.ID,
				RoleID: role.ID,
			}

			if err := tx.Create(userRole).Error; err != nil {
				log.Warn("分配管理员角色失败", "error", err)
				return nil
			}

			log.Info("管理员用户创建成功",
				"username", username,
				"nickname", nickname,
				"id", sysUser.ID)
			return nil
		},
		Down: func(tx *gorm.DB) error {
			// 删除管理员用户
			username := "admin"
			if migrator.config != nil && migrator.config.Admin.Username != "" {
				username = migrator.config.Admin.Username
			}

			return tx.Where("username = ?", username).Delete(&model.SysUser{}).Error
		},
	})
}

// RegisterAllMigrations 注册所有迁移
func RegisterAllMigrations(migrator *Migrator) {
	// 按顺序注册所有迁移
	RegisterCoreUserMigrations(migrator)
	RegisterRoleMigrations(migrator)
	RegisterPermissionMigrations(migrator)
	RegisterMenuMigrations(migrator)
	RegisterOperationLogMigrations(migrator)
	RegisterInitialDataMigrations(migrator) // 添加初始数据迁移

	// 添加自定义业务迁移...
}
