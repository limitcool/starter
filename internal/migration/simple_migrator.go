package migration

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/gorm"
)

// SimpleMigration 迁移记录
type SimpleMigration struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"uniqueIndex;size:255;not null;comment:迁移名称"`
	CreatedAt time.Time `gorm:"comment:执行时间"`
	Batch     int       `gorm:"comment:批次"`
}

// TableName 指定迁移表名
func (SimpleMigration) TableName() string {
	return "migrations"
}

// SimpleMigrator 简化版迁移管理器
type SimpleMigrator struct {
	db     *gorm.DB
	config *configs.Config
}

// NewSimpleMigrator 创建简化版迁移管理器
func NewSimpleMigrator(db *gorm.DB, config *configs.Config) *SimpleMigrator {
	return &SimpleMigrator{
		db:     db,
		config: config,
	}
}

// Initialize 初始化迁移系统
func (m *SimpleMigrator) Initialize() error {
	// 创建迁移表
	if err := m.db.AutoMigrate(&SimpleMigration{}); err != nil {
		return fmt.Errorf("创建迁移表失败: %w", err)
	}

	logger.Info("迁移系统初始化完成")
	return nil
}

// Migrate 执行所有未运行的迁移
func (m *SimpleMigrator) Migrate() error {
	// 先检查迁移表是否存在
	if err := m.Initialize(); err != nil {
		return err
	}

	// 获取已运行的迁移
	var ranMigrations []SimpleMigration
	if err := m.db.Find(&ranMigrations).Error; err != nil {
		return fmt.Errorf("获取已运行迁移失败: %w", err)
	}

	ranMigrationNames := make(map[string]struct{})
	for _, migration := range ranMigrations {
		ranMigrationNames[migration.Name] = struct{}{}
	}

	// 获取最新批次号
	var lastBatch int
	m.db.Model(&SimpleMigration{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch)
	currentBatch := lastBatch + 1

	// 获取所有迁移
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("加载迁移失败: %w", err)
	}

	// 执行未运行的迁移
	for _, migration := range migrations {
		if _, ok := ranMigrationNames[migration.Name]; !ok {
			logger.Info("执行迁移", "name", migration.Name)

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
			record := SimpleMigration{
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

			logger.Info("迁移完成", "name", migration.Name)
		}
	}

	return nil
}

// Rollback 回滚最后一批迁移
func (m *SimpleMigrator) Rollback() error {
	// 获取最后一批迁移
	var lastMigrations []SimpleMigration
	var lastBatch int

	if err := m.db.Model(&SimpleMigration{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch).Error; err != nil {
		return fmt.Errorf("获取最后批次失败: %w", err)
	}

	if lastBatch == 0 {
		return errors.New("没有可回滚的迁移")
	}

	if err := m.db.Where("batch = ?", lastBatch).Order("id DESC").Find(&lastMigrations).Error; err != nil {
		return fmt.Errorf("获取最后批次迁移失败: %w", err)
	}

	// 获取所有迁移
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("加载迁移失败: %w", err)
	}

	// 构建名称到迁移的映射
	migrationsMap := make(map[string]*SimpleMigrationItem)
	for _, migration := range migrations {
		migrationsMap[migration.Name] = migration
	}

	// 回滚迁移
	for _, migration := range lastMigrations {
		migrationItem, ok := migrationsMap[migration.Name]
		if !ok || migrationItem.Down == nil {
			logger.Warn("没有找到回滚函数", "name", migration.Name)
			continue
		}

		logger.Info("回滚迁移", "name", migration.Name)

		// 开始事务
		tx := m.db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("开始事务失败: %w", tx.Error)
		}

		// 执行回滚
		if err := migrationItem.Down(tx); err != nil {
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

		logger.Info("回滚完成", "name", migration.Name)
	}

	return nil
}

// Reset 重置所有迁移
func (m *SimpleMigrator) Reset() error {
	// 获取所有迁移，按反向顺序
	var allMigrations []SimpleMigration
	if err := m.db.Order("id DESC").Find(&allMigrations).Error; err != nil {
		return fmt.Errorf("获取所有迁移失败: %w", err)
	}

	if len(allMigrations) == 0 {
		return errors.New("没有可重置的迁移")
	}

	// 获取所有迁移
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("加载迁移失败: %w", err)
	}

	// 构建名称到迁移的映射
	migrationsMap := make(map[string]*SimpleMigrationItem)
	for _, migration := range migrations {
		migrationsMap[migration.Name] = migration
	}

	// 回滚所有迁移
	for _, migration := range allMigrations {
		migrationItem, ok := migrationsMap[migration.Name]
		if !ok || migrationItem.Down == nil {
			logger.Warn("没有找到回滚函数", "name", migration.Name)
			continue
		}

		logger.Info("重置迁移", "name", migration.Name)

		// 开始事务
		tx := m.db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("开始事务失败: %w", tx.Error)
		}

		// 执行回滚
		if err := migrationItem.Down(tx); err != nil {
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

		logger.Info("重置完成", "name", migration.Name)
	}

	return nil
}

// Status 获取迁移状态
func (m *SimpleMigrator) Status() ([]map[string]any, error) {
	// 获取所有已运行的迁移
	var ranMigrations []SimpleMigration
	if err := m.db.Find(&ranMigrations).Error; err != nil {
		return nil, fmt.Errorf("获取已运行迁移失败: %w", err)
	}

	ranMigrationMap := make(map[string]SimpleMigration)
	for _, migration := range ranMigrations {
		ranMigrationMap[migration.Name] = migration
	}

	// 获取所有迁移
	migrations, err := m.loadMigrations()
	if err != nil {
		return nil, fmt.Errorf("加载迁移失败: %w", err)
	}

	// 构建状态
	var status []map[string]any
	for _, migration := range migrations {
		ran, exists := ranMigrationMap[migration.Name]
		var batch int
		var timestamp time.Time
		if exists {
			batch = ran.Batch
			timestamp = ran.CreatedAt
		}
		status = append(status, map[string]any{
			"name":      migration.Name,
			"ran":       exists,
			"batch":     batch,
			"timestamp": timestamp,
		})
	}

	return status, nil
}

// SimpleMigrationItem 表示单个迁移项
type SimpleMigrationItem struct {
	Name string                // 迁移名称，如 "create_users_table"
	Up   func(*gorm.DB) error // 向上迁移函数
	Down func(*gorm.DB) error // 向下迁移函数
}

// loadMigrations 加载所有迁移
func (m *SimpleMigrator) loadMigrations() ([]*SimpleMigrationItem, error) {
	// 这里我们将返回所有注册的迁移
	// 在实际实现中，你可以从文件系统或其他地方加载迁移
	migrations := make([]*SimpleMigrationItem, 0)

	// 添加所有迁移
	for _, migration := range GetAllMigrations() {
		migrations = append(migrations, migration)
	}

	// 按名称排序
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	return migrations, nil
}

// InitializeSimpleMigrator 初始化简化版迁移器
func InitializeSimpleMigrator(db *gorm.DB, config *configs.Config) (*SimpleMigrator, error) {
	migrator := NewSimpleMigrator(db, config)
	err := migrator.Initialize()
	return migrator, err
}

// 迁移注册表
var migrationRegistry = make([]*SimpleMigrationItem, 0)

// RegisterMigration 注册迁移
func RegisterMigration(name string, up, down func(*gorm.DB) error) {
	migrationRegistry = append(migrationRegistry, &SimpleMigrationItem{
		Name: name,
		Up:   up,
		Down: down,
	})
}

// GetAllMigrations 获取所有注册的迁移
func GetAllMigrations() []*SimpleMigrationItem {
	return migrationRegistry
}

// ScanMigrationFiles 扫描迁移文件目录
func ScanMigrationFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理Go文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_migration.go") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 按文件名排序
	sort.Strings(files)

	return files, nil
}
