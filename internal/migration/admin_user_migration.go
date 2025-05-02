package migration

import (
	"fmt"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// RegisterAdminUserMigrations 注册管理员用户相关迁移
func RegisterAdminUserMigrations(migrator *Migrator) {
	// 注册管理员用户迁移
	registerAdminUserMigration(migrator)
}

// registerAdminUserMigration 注册管理员用户迁移
func registerAdminUserMigration(migrator *Migrator) {
	// 创建管理员用户表
	migrator.Register(&MigrationEntry{
		Version: "202504080301",
		Name:    "create_admin_user_table",
		Up: func(tx *gorm.DB) error {
			// 检查表是否已存在
			if tx.Migrator().HasTable("admin_user") {
				logger.Info("admin_user表已存在，跳过创建")
				return nil
			}

			// 创建admin_user表
			return tx.AutoMigrate(&model.AdminUser{})
		},
		Down: func(tx *gorm.DB) error {
			// 删除admin_user表
			return tx.Migrator().DropTable("admin_user")
		},
	})

	// 创建admin_user_role关联表
	migrator.Register(&MigrationEntry{
		Version: "202504080302",
		Name:    "create_admin_user_role_table",
		Up: func(tx *gorm.DB) error {
			// 检查表是否已存在
			if tx.Migrator().HasTable("admin_user_role") {
				logger.Info("admin_user_role表已存在，跳过创建")
				return nil
			}

			// 创建admin_user_role关联表
			type AdminUserRole struct {
				model.BaseModel
				AdminUserID int64 `gorm:"type:bigint;not null;index;comment:管理员用户ID"`
				RoleID      uint  `gorm:"not null;index;comment:角色ID"`
			}

			return tx.AutoMigrate(&AdminUserRole{})
		},
		Down: func(tx *gorm.DB) error {
			// 删除admin_user_role表
			return tx.Migrator().DropTable("admin_user_role")
		},
	})

	// 初始化分离模式下的管理员用户
	migrator.Register(&MigrationEntry{
		Version: "202504080303",
		Name:    "init_separate_admin_user",
		Up: func(tx *gorm.DB) error {
			// 获取配置
			cfg := migrator.config
			if cfg == nil {
				logger.Warn("配置未初始化，使用默认管理员账号")
				cfg = &configs.Config{
					Admin: configs.Admin{
						Username: "admin",
						Password: "123456",
						Nickname: "超级管理员",
					},
				}
			}

			// 获取管理员信息
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

			// 已存在则跳过
			if count > 0 {
				logger.Info("分离模式管理员用户已存在，跳过创建")
				return nil
			}

			// 创建管理员用户
			hashedPassword, err := crypto.HashPassword(password)
			if err != nil {
				return fmt.Errorf("密码加密失败: %w", err)
			}

			// 创建管理员用户
			adminUser := &model.AdminUser{
				Username: username,
				Password: hashedPassword,
				Nickname: nickname,
				Enabled:  true,
				Remark:   "系统初始化创建",
				Email:    "",
				Mobile:   "",
			}

			logger.Info("准备创建分离模式管理员用户",
				"username", adminUser.Username,
				"nickname", adminUser.Nickname)

			if err := tx.Create(adminUser).Error; err != nil {
				return fmt.Errorf("创建分离模式管理员账号失败: %w", err)
			}

			// 获取管理员角色
			var adminRole model.Role
			if err := tx.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
				logger.Warn("获取管理员角色失败，跳过角色分配", "error", err)
				return nil
			}

			// 创建用户角色关联
			adminUserRole := struct {
				AdminUserID int64 `gorm:"column:admin_user_id"`
				RoleID      uint  `gorm:"column:role_id"`
			}{
				AdminUserID: adminUser.ID,
				RoleID:      adminRole.ID,
			}

			if err := tx.Table("admin_user_role").Create(&adminUserRole).Error; err != nil {
				logger.Warn("分配管理员角色失败", "error", err)
				return nil
			}

			logger.Info("分离模式管理员用户创建成功",
				"username", username,
				"nickname", nickname,
				"id", adminUser.ID)
			return nil
		},
		Down: func(tx *gorm.DB) error {
			// 删除分离模式管理员用户
			username := "admin"
			if migrator.config != nil && migrator.config.Admin.Username != "" {
				username = migrator.config.Admin.Username
			}

			return tx.Where("username = ?", username).Delete(&model.AdminUser{}).Error
		},
	})
}
