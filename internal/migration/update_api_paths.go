package migration

import (
	"context"
	"strings"

	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// RegisterUpdateAPIPathsMigrations 注册更新API路径的迁移
func RegisterUpdateAPIPathsMigrations(migrator *Migrator) {
	// 更新API路径
	migrator.Register(&MigrationEntry{
		Version: "202504210001",
		Name:    "update_api_paths",
		Up: func(tx *gorm.DB) error {
			// 查询所有包含 /api/v1/admin/ 的API
			var apis []model.API
			if err := tx.Where("path LIKE ?", "%/api/v1/admin/%").Find(&apis).Error; err != nil {
				return err
			}

			// 如果没有找到API，直接返回
			if len(apis) == 0 {
				logger.InfoContext(context.Background(), "没有找到需要更新的API路径")
				return nil
			}

			// 更新API路径
			for _, api := range apis {
				// 替换路径中的 /api/v1/admin/ 为 /api/v1/admin-api/
				newPath := strings.Replace(api.Path, "/api/v1/admin/", "/api/v1/admin-api/", 1)

				// 更新API路径
				if err := tx.Model(&model.API{}).Where("id = ?", api.ID).Update("path", newPath).Error; err != nil {
					return err
				}

				// 更新权限记录中的code
				oldCode := api.Path + ":" + api.Method
				newCode := newPath + ":" + api.Method

				if err := tx.Model(&model.Permission{}).Where("code = ?", oldCode).Update("code", newCode).Error; err != nil {
					return err
				}

				logger.InfoContext(context.Background(), "更新API路径", "old_path", api.Path, "new_path", newPath)
			}

			// 更新Casbin规则中的路径
			var rules []CasbinRule
			if err := tx.Where("v1 LIKE ?", "%/api/v1/admin/%").Find(&rules).Error; err != nil {
				return err
			}

			for _, rule := range rules {
				// 替换路径中的 /api/v1/admin/ 为 /api/v1/admin-api/
				newV1 := strings.Replace(rule.V1, "/api/v1/admin/", "/api/v1/admin-api/", 1)

				// 更新规则
				if err := tx.Model(&CasbinRule{}).Where("id = ?", rule.ID).Update("v1", newV1).Error; err != nil {
					return err
				}

				logger.InfoContext(context.Background(), "更新Casbin规则", "old_v1", rule.V1, "new_v1", newV1)
			}

			return nil
		},
		Down: func(tx *gorm.DB) error {
			// 查询所有包含 /api/v1/admin-api/ 的API
			var apis []model.API
			if err := tx.Where("path LIKE ?", "%/api/v1/admin-api/%").Find(&apis).Error; err != nil {
				return err
			}

			// 如果没有找到API，直接返回
			if len(apis) == 0 {
				logger.InfoContext(context.Background(), "没有找到需要回滚的API路径")
				return nil
			}

			// 回滚API路径
			for _, api := range apis {
				// 替换路径中的 /api/v1/admin-api/ 为 /api/v1/admin/
				oldPath := api.Path
				newPath := strings.Replace(api.Path, "/api/v1/admin-api/", "/api/v1/admin/", 1)

				// 更新API路径
				if err := tx.Model(&model.API{}).Where("id = ?", api.ID).Update("path", newPath).Error; err != nil {
					return err
				}

				// 更新权限记录中的code
				oldCode := oldPath + ":" + api.Method
				newCode := newPath + ":" + api.Method

				if err := tx.Model(&model.Permission{}).Where("code = ?", oldCode).Update("code", newCode).Error; err != nil {
					return err
				}

				logger.Info("回滚API路径", "old_path", oldPath, "new_path", newPath)
			}

			// 回滚Casbin规则中的路径
			var rules []CasbinRule
			if err := tx.Where("v1 LIKE ?", "%/api/v1/admin-api/%").Find(&rules).Error; err != nil {
				return err
			}

			for _, rule := range rules {
				// 替换路径中的 /api/v1/admin-api/ 为 /api/v1/admin/
				newV1 := strings.Replace(rule.V1, "/api/v1/admin-api/", "/api/v1/admin/", 1)

				// 更新规则
				if err := tx.Model(&CasbinRule{}).Where("id = ?", rule.ID).Update("v1", newV1).Error; err != nil {
					return err
				}

				logger.Info("回滚Casbin规则", "old_v1", rule.V1, "new_v1", newV1)
			}

			return nil
		},
	})
}
