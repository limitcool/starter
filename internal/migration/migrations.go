package migration

import (
	"fmt"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// RegisterMigrations 注册数据库迁移
func RegisterMigrations(migrator *Migrator) {
	// 添加文件表迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080000",
		Name:    "create_file_table",
		Up: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.File{})
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("file")
		},
	})

	// 添加普通用户表迁移
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

	// 添加初始管理员用户迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080010",
		Name:    "init_admin_user",
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
			if err := tx.Model(&model.User{}).Where("username = ? AND is_admin = ?", username, true).Count(&count).Error; err != nil {
				return err
			}

			// 已存在则不重复创建
			if count > 0 {
				logger.Info("管理员用户已存在，跳过创建")
				return nil
			}

			// 创建管理员用户
			hashedPassword, err := crypto.HashPassword(password)
			if err != nil {
				return fmt.Errorf("密码加密失败: %w", err)
			}

			// 创建管理员用户（在普通用户表中，使用is_admin字段标识）
			adminUser := &model.User{
				Username: username,
				Password: hashedPassword,
				Nickname: nickname,
				Email:    "admin@example.com",
				Enabled:  true,
				IsAdmin:  true, // 标记为管理员
			}

			logger.Info("准备创建管理员用户",
				"username", adminUser.Username,
				"nickname", adminUser.Nickname)

			if err := tx.Create(adminUser).Error; err != nil {
				return fmt.Errorf("创建管理员账号失败: %w", err)
			}

			logger.Info("管理员用户创建成功",
				"username", username,
				"nickname", nickname,
				"id", adminUser.ID)
			return nil
		},
		Down: func(tx *gorm.DB) error {
			// 删除管理员用户
			username := "admin"
			if migrator.config != nil && migrator.config.Admin.Username != "" {
				username = migrator.config.Admin.Username
			}

			return tx.Where("username = ? AND is_admin = ?", username, true).Delete(&model.User{}).Error
		},
	})
}
