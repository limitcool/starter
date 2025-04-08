package migration

import (
	"github.com/limitcool/starter/internal/model"
	"gorm.io/gorm"
)

// 此文件用于定义具体的迁移实例
// 您可以在这里添加新的迁移，按照时间顺序排列

// RegisterCoreUserMigrations 注册用户相关迁移
func RegisterCoreUserMigrations(migrator *Migrator) {
	// 用户表迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080001",
		Name:    "create_users_table",
		Up: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.User{})
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("sys_user")
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
}

// RegisterPermissionMigrations 注册权限相关迁移
func RegisterPermissionMigrations(migrator *Migrator) {
	// 权限表迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080003",
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
		Version: "202504080004",
		Name:    "create_menus_table",
		Up: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Menu{})
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("sys_menu")
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

	// 添加自定义业务迁移...
}
