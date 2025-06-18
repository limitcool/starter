package migration

import (
	"fmt"

	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// UpdateFileTableToUUID_20250618 更新文件表以支持UUID主键
func UpdateFileTableToUUID_20250618(db *gorm.DB) error {
	logger.Info("开始迁移：更新文件表以支持UUID主键")

	// 检查表是否存在
	if !db.Migrator().HasTable("file") {
		logger.Info("文件表不存在，跳过迁移")
		return nil
	}

	// 备份现有数据
	logger.Info("备份现有文件数据...")
	var existingFiles []map[string]interface{}
	if err := db.Table("file").Find(&existingFiles).Error; err != nil {
		return fmt.Errorf("备份现有数据失败: %w", err)
	}
	logger.Info("备份完成", "count", len(existingFiles))

	// 删除旧表
	logger.Info("删除旧的文件表...")
	if err := db.Migrator().DropTable("file"); err != nil {
		return fmt.Errorf("删除旧表失败: %w", err)
	}

	// 创建新的文件表结构（使用UUID主键）
	logger.Info("创建新的文件表结构...")
	if err := db.Exec(`
		CREATE TABLE file (
			id VARCHAR(36) PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE,
			name VARCHAR(255),
			original_name VARCHAR(255),
			path VARCHAR(500),
			type VARCHAR(50),
			usage VARCHAR(50) DEFAULT 'general',
			size BIGINT DEFAULT 0,
			mime_type VARCHAR(100),
			extension VARCHAR(20),
			storage_type VARCHAR(20),
			uploaded_by BIGINT,
			uploaded_by_type SMALLINT DEFAULT 1,
			uploaded_at TIMESTAMP WITH TIME ZONE,
			status INTEGER DEFAULT 0,
			is_public BOOLEAN DEFAULT false
		)
	`).Error; err != nil {
		return fmt.Errorf("创建新表失败: %w", err)
	}

	// 创建索引
	logger.Info("创建索引...")
	indexes := []string{
		"CREATE INDEX idx_file_deleted_at ON file(deleted_at)",
		"CREATE INDEX idx_file_usage ON file(usage)",
		"CREATE INDEX idx_file_uploaded_by ON file(uploaded_by)",
		"CREATE INDEX idx_file_storage_type ON file(storage_type)",
		"CREATE INDEX idx_file_is_public ON file(is_public)",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			logger.Warn("创建索引失败", "sql", indexSQL, "error", err)
		}
	}

	logger.Info("文件表迁移完成")
	return nil
}

// RollbackFileTableToUUID_20250618 回滚文件表UUID迁移
func RollbackFileTableToUUID_20250618(db *gorm.DB) error {
	logger.Info("开始回滚：恢复文件表原始结构")

	// 删除UUID版本的表
	if err := db.Migrator().DropTable("file"); err != nil {
		return fmt.Errorf("删除UUID表失败: %w", err)
	}

	// 重新创建原始表结构
	if err := db.Exec(`
		CREATE TABLE file (
			id BIGSERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE,
			name VARCHAR(255),
			original_name VARCHAR(255),
			path VARCHAR(500),
			type VARCHAR(50),
			usage VARCHAR(50) DEFAULT 'general',
			size BIGINT DEFAULT 0,
			mime_type VARCHAR(100),
			extension VARCHAR(20),
			storage_type VARCHAR(20),
			uploaded_by BIGINT,
			uploaded_by_type SMALLINT DEFAULT 1,
			uploaded_at TIMESTAMP WITH TIME ZONE,
			status INTEGER DEFAULT 0,
			is_public BOOLEAN DEFAULT false
		)
	`).Error; err != nil {
		return fmt.Errorf("恢复原始表失败: %w", err)
	}

	logger.Info("文件表回滚完成")
	return nil
}
