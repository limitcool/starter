package migration

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// Migration 迁移记录
type Migration struct {
	ID        uint      `gorm:"primarykey"`
	Version   string    `gorm:"uniqueIndex;size:50;not null;comment:版本号"`
	Name      string    `gorm:"size:100;not null;comment:迁移名称"`
	CreatedAt time.Time `gorm:"comment:执行时间"`
	Batch     int       `gorm:"comment:批次"`
}

// TableName 指定迁移表名
func (Migration) TableName() string {
	return "sys_migrations"
}

// MigrationFunc 定义迁移函数类型
type MigrationFunc func(*gorm.DB) error

// MigrationEntry 表示单个迁移项
type MigrationEntry struct {
	Version string        // 版本号，如 202504080001
	Name    string        // 迁移名称，如 "create_users_table"
	Up      MigrationFunc // 向上迁移函数
	Down    MigrationFunc // 向下迁移函数
}

// Migrator 迁移管理器
type Migrator struct {
	db         *gorm.DB
	migrations []*MigrationEntry
	config     *configs.Config
}

// NewMigrator 创建迁移管理器
func NewMigrator(db *gorm.DB, config *configs.Config) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]*MigrationEntry, 0),
		config:     config,
	}
}

// Register 注册迁移项
func (m *Migrator) Register(migration *MigrationEntry) {
	m.migrations = append(m.migrations, migration)
}

// Initialize 初始化迁移系统
func (m *Migrator) Initialize() error {
	// 创建迁移表
	if err := m.db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("创建迁移表失败: %w", err)
	}

	logger.Info("迁移系统初始化完成")
	return nil
}

// Migrate 执行所有未运行的迁移
func (m *Migrator) Migrate() error {
	// 先检查迁移表是否存在
	if err := m.Initialize(); err != nil {
		return err
	}

	// 获取已运行的迁移
	var ranMigrations []Migration
	if err := m.db.Find(&ranMigrations).Error; err != nil {
		return fmt.Errorf("获取已运行迁移失败: %w", err)
	}

	ranVersions := make(map[string]struct{})
	for _, migration := range ranMigrations {
		ranVersions[migration.Version] = struct{}{}
	}

	// 排序迁移，确保按版本号顺序执行
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	// 获取最新批次号
	var lastBatch int
	m.db.Model(&Migration{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch)
	currentBatch := lastBatch + 1

	// 执行未运行的迁移
	for _, migration := range m.migrations {
		if _, ok := ranVersions[migration.Version]; !ok {
			logger.Info("执行迁移", "version", migration.Version, "name", migration.Name)

			// 开始事务
			tx := m.db.Begin()
			if tx.Error != nil {
				return fmt.Errorf("开始事务失败: %w", tx.Error)
			}

			// 执行迁移
			if err := migration.Up(tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("迁移失败 (%s): %w", migration.Name, err)
			}

			// 记录迁移
			record := Migration{
				Version:   migration.Version,
				Name:      migration.Name,
				CreatedAt: time.Now(),
				Batch:     currentBatch,
			}
			if err := tx.Create(&record).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("记录迁移失败: %w", err)
			}

			// 提交事务
			if err := tx.Commit().Error; err != nil {
				return fmt.Errorf("提交事务失败: %w", err)
			}

			logger.Info("迁移完成", "version", migration.Version, "name", migration.Name)
		}
	}

	return nil
}

// Rollback 回滚最后一批迁移
func (m *Migrator) Rollback() error {
	// 获取最后一批迁移
	var lastMigrations []Migration
	var lastBatch int

	if err := m.db.Model(&Migration{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch).Error; err != nil {
		return fmt.Errorf("获取最后批次失败: %w", err)
	}

	if lastBatch == 0 {
		return errors.New("没有可回滚的迁移")
	}

	if err := m.db.Where("batch = ?", lastBatch).Order("id DESC").Find(&lastMigrations).Error; err != nil {
		return fmt.Errorf("获取最后批次迁移失败: %w", err)
	}

	// 构建版本到迁移的映射
	migrationsMap := make(map[string]*MigrationEntry)
	for _, migration := range m.migrations {
		migrationsMap[migration.Version] = migration
	}

	// 回滚迁移
	for _, migration := range lastMigrations {
		migrationEntry, ok := migrationsMap[migration.Version]
		if !ok || migrationEntry.Down == nil {
			logger.Warn("没有找到回滚函数", "version", migration.Version, "name", migration.Name)
			continue
		}

		logger.Info("回滚迁移", "version", migration.Version, "name", migration.Name)

		// 开始事务
		tx := m.db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("开始事务失败: %w", tx.Error)
		}

		// 执行回滚
		if err := migrationEntry.Down(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("回滚失败 (%s): %w", migration.Name, err)
		}

		// 删除迁移记录
		if err := tx.Delete(&migration).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("删除迁移记录失败: %w", err)
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("提交事务失败: %w", err)
		}

		logger.Info("回滚完成", "version", migration.Version, "name", migration.Name)
	}

	return nil
}

// Reset 重置所有迁移
func (m *Migrator) Reset() error {
	// 获取所有迁移，按反向顺序
	var allMigrations []Migration
	if err := m.db.Order("id DESC").Find(&allMigrations).Error; err != nil {
		return fmt.Errorf("获取所有迁移失败: %w", err)
	}

	if len(allMigrations) == 0 {
		return errors.New("没有可重置的迁移")
	}

	// 构建版本到迁移的映射
	migrationsMap := make(map[string]*MigrationEntry)
	for _, migration := range m.migrations {
		migrationsMap[migration.Version] = migration
	}

	// 回滚所有迁移
	for _, migration := range allMigrations {
		migrationEntry, ok := migrationsMap[migration.Version]
		if !ok || migrationEntry.Down == nil {
			logger.Warn("没有找到回滚函数", "version", migration.Version, "name", migration.Name)
			continue
		}

		logger.Info("重置迁移", "version", migration.Version, "name", migration.Name)

		// 开始事务
		tx := m.db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("开始事务失败: %w", tx.Error)
		}

		// 执行回滚
		if err := migrationEntry.Down(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("重置失败 (%s): %w", migration.Name, err)
		}

		// 删除迁移记录
		if err := tx.Delete(&migration).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("删除迁移记录失败: %w", err)
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("提交事务失败: %w", err)
		}

		logger.Info("重置完成", "version", migration.Version, "name", migration.Name)
	}

	return nil
}

// Status 获取迁移状态
func (m *Migrator) Status() ([]map[string]any, error) {
	// 获取所有已运行的迁移
	var ranMigrations []Migration
	if err := m.db.Find(&ranMigrations).Error; err != nil {
		return nil, fmt.Errorf("获取已运行迁移失败: %w", err)
	}

	ranVersions := make(map[string]Migration)
	for _, migration := range ranMigrations {
		ranVersions[migration.Version] = migration
	}

	// 排序迁移
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	// 构建状态
	var status []map[string]any
	for _, migration := range m.migrations {
		ran, exists := ranVersions[migration.Version]
		var batch int
		var timestamp time.Time
		if exists {
			batch = ran.Batch
			timestamp = ran.CreatedAt
		}
		status = append(status, map[string]any{
			"version":   migration.Version,
			"name":      migration.Name,
			"ran":       exists,
			"batch":     batch,
			"timestamp": timestamp,
		})
	}

	return status, nil
}

// InitializeMigrator 初始化迁移器并注册所有迁移
func InitializeMigrator(db *gorm.DB, config *configs.Config) (*Migrator, error) {
	migrator := NewMigrator(db, config)
	RegisterAllMigrations(migrator) // 注册所有迁移

	err := migrator.Initialize()
	return migrator, err
}
