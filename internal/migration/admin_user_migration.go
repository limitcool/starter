package migration

import (
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// RegisterAdminUserMigrations 注册管理员用户相关迁移
func RegisterAdminUserMigrations(migrator *Migrator) {
	// 管理员用户表迁移
	migrator.Register(&MigrationEntry{
		Version: "202504080101",
		Name:    "create_admin_users_table",
		Up: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.AdminUser{})
		},
		Down: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("admin_user")
		},
	})

	// 从sys_user表迁移数据到admin_user表
	migrator.Register(&MigrationEntry{
		Version: "202504080102",
		Name:    "migrate_sys_users_to_admin_users",
		Up: func(tx *gorm.DB) error {
			// 检查sys_user表是否存在
			if !tx.Migrator().HasTable("sys_user") {
				logger.Info("sys_user表不存在，跳过数据迁移")
				return nil
			}

			// 检查admin_user表是否已有数据
			var count int64
			if err := tx.Model(&model.AdminUser{}).Count(&count).Error; err != nil {
				return err
			}

			// 如果admin_user表已有数据，跳过迁移
			if count > 0 {
				logger.Info("admin_user表已有数据，跳过数据迁移")
				return nil
			}

			// 执行数据迁移
			// 注意：这里使用原始SQL，因为我们需要直接操作表，而不是通过模型
			sql := `INSERT INTO admin_user 
					SELECT * FROM sys_user`

			if err := tx.Exec(sql).Error; err != nil {
				return err
			}

			// 更新关联表
			if tx.Migrator().HasTable("sys_user_role") {
				// 创建admin_user_role表
				if err := tx.AutoMigrate(&model.UserRole{}); err != nil {
					return err
				}

				// 修改UserRole的表名
				if err := tx.Exec("ALTER TABLE sys_user_role RENAME TO admin_user_role").Error; err != nil {
					logger.Warn("重命名sys_user_role表失败，可能已经重命名", "error", err)
				}
			}

			logger.Info("从sys_user表迁移数据到admin_user表成功")
			return nil
		},
		Down: func(tx *gorm.DB) error {
			// 这里不提供回滚操作，因为可能会导致数据丢失
			logger.Warn("不支持从admin_user表回滚到sys_user表")
			return nil
		},
	})
}
