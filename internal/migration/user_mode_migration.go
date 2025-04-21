package migration

import (
	"fmt"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// RegisterUserModeMigrations 注册用户模式相关迁移
func RegisterUserModeMigrations(migrator *Migrator) {
	// 更新用户表，添加is_admin字段
	migrator.Register(&MigrationEntry{
		Version: "202504080201",
		Name:    "add_is_admin_to_user_table",
		Up: func(tx *gorm.DB) error {
			// 检查字段是否已存在
			if tx.Migrator().HasColumn(&model.User{}, "is_admin") {
				logger.Info("is_admin字段已存在，跳过添加")
				return nil
			}

			// 添加is_admin字段
			return tx.Migrator().AddColumn(&model.User{}, "is_admin")
		},
		Down: func(tx *gorm.DB) error {
			// 删除is_admin字段
			return tx.Migrator().DropColumn(&model.User{}, "is_admin")
		},
	})

	// 更新用户表，添加角色关联
	migrator.Register(&MigrationEntry{
		Version: "202504080202",
		Name:    "create_user_role_table",
		Up: func(tx *gorm.DB) error {
			// 检查表是否已存在
			if tx.Migrator().HasTable("user_role") {
				logger.Info("user_role表已存在，跳过创建")
				return nil
			}

			// 创建user_role关联表
			type UserRole struct {
				model.BaseModel
				UserID int64 `gorm:"type:bigint;not null;index;comment:用户ID"`
				RoleID uint  `gorm:"not null;index;comment:角色ID"`
			}

			return tx.AutoMigrate(&UserRole{})
		},
		Down: func(tx *gorm.DB) error {
			// 删除user_role表
			return tx.Migrator().DropTable("user_role")
		},
	})

	// 初始化合并模式下的管理员用户
	migrator.Register(&MigrationEntry{
		Version: "202504080203",
		Name:    "init_unified_admin_user",
		Up: func(tx *gorm.DB) error {
			// 获取配置
			cfg := migrator.config
			if cfg == nil {
				logger.Warn("配置未初始化，跳过创建合并模式管理员")
				return nil
			}

			// 检查是否是合并模式
			userMode := enum.GetUserMode(cfg.Admin.UserMode)
			if userMode != enum.UserModeSimple {
				logger.Info("非简单模式，跳过创建简单模式管理员")
				return nil
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
			if err := tx.Model(&model.User{}).Where("username = ? AND is_admin = ?", username, true).Count(&count).Error; err != nil {
				return err
			}

			// 已存在则跳过
			if count > 0 {
				logger.Info("简单模式管理员用户已存在，跳过创建")
				return nil
			}

			// 创建管理员用户
			hashedPassword, err := crypto.HashPassword(password)
			if err != nil {
				return fmt.Errorf("密码加密失败: %w", err)
			}

			// 创建管理员用户
			adminUser := &model.User{
				Username: username,
				Password: hashedPassword,
				Nickname: nickname,
				Enabled:  true,
				Remark:   "系统初始化创建",
				Email:    "",
				Mobile:   "",
				IsAdmin:  true, // 设置为管理员
			}

			logger.Info("准备创建简单模式管理员用户",
				"username", adminUser.Username,
				"nickname", adminUser.Nickname)

			if err := tx.Create(adminUser).Error; err != nil {
				return fmt.Errorf("创建简单模式管理员账号失败: %w", err)
			}

			logger.Info("简单模式管理员用户创建成功",
				"username", username,
				"nickname", nickname,
				"id", adminUser.ID)
			return nil
		},
		Down: func(tx *gorm.DB) error {
			// 删除简单模式管理员用户
			username := "admin"
			if migrator.config != nil && migrator.config.Admin.Username != "" {
				username = migrator.config.Admin.Username
			}

			return tx.Where("username = ? AND is_admin = ?", username, true).Delete(&model.User{}).Error
		},
	})
}
