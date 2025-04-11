package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/internal/migration"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"github.com/spf13/cobra"
)

// migrateCmd 表示migrate子命令
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Execute database migrations",
	Long: `Execute database migrations to create or update database table structures.

This command automatically performs database migration operations based on defined model structures,
suitable for development environments and initializing production environments.
Make sure to back up your database before using it in a production environment.`,
	Run: runMigration,
}

// migrateRollbackCmd 表示migrate:rollback子命令
var migrateRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback the last batch of database migrations",
	Long:  `Rollback the last batch of database migration operations, restoring the database to its previous state.`,
	Run:   runMigrationRollback,
}

// migrateStatusCmd 表示migrate:status子命令
var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display database migration status",
	Long:  `Display all migrations and their execution status.`,
	Run:   runMigrationStatus,
}

// migrateResetCmd 表示migrate:reset子命令
var migrateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all database migrations",
	Long:  `Reset all executed database migrations, restoring the database to its initial state.`,
	Run:   runMigrationReset,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateRollbackCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migrateResetCmd)

	// 添加迁移特定的标志
	migrateCmd.PersistentFlags().BoolP("fresh", "f", false, "Clear the database before migration (dangerous operation)")
}

// runMigration 执行数据库迁移
func runMigration(cmd *cobra.Command, args []string) {
	// 加载配置
	cfg := InitConfig(cmd, args)

	// 设置日志
	InitLogger(cfg)

	// 检查数据库是否启用
	if !cfg.Database.Enabled {
		log.Fatal("Database not enabled, please enable it in the configuration file")
	}

	log.Info("Starting database migration process")

	// 初始化SQL组件
	sqlComponent := sqldb.NewComponent(cfg)
	if err := sqlComponent.Initialize(); err != nil {
		log.Error("Failed to initialize SQL component", "error", err)
		os.Exit(1)
	}
	defer sqlComponent.Cleanup()

	// 初始化迁移系统
	if err := migration.InitializeMigrator(sqlComponent.DB(), cfg); err != nil {
		log.Error("Failed to initialize migration system", "error", err)
		os.Exit(1)
	}

	// 检查是否需要重置数据库
	fresh, _ := cmd.Flags().GetBool("fresh")
	if fresh {
		log.Warn("Preparing to clear and rebuild the database...")
		if err := migration.GlobalMigrator.Reset(); err != nil {
			if err.Error() != "No migrations to reset" {
				log.Error("Failed to reset database", "error", err)
				os.Exit(1)
			}
		}
	}

	// 执行迁移
	if err := migration.GlobalMigrator.Migrate(); err != nil {
		log.Error("Database migration failed", "error", err)
		os.Exit(1)
	}

	log.Info("Database migration completed successfully")
}

// runMigrationRollback 回滚数据库迁移
func runMigrationRollback(cmd *cobra.Command, args []string) {
	// 加载配置
	cfg := InitConfig(cmd, args)

	// 设置日志
	InitLogger(cfg)

	// 检查数据库是否启用
	if !cfg.Database.Enabled {
		log.Fatal("Database not enabled, please enable it in the configuration file")
	}

	log.Info("Starting database migration rollback process")

	// 初始化SQL组件
	sqlComponent := sqldb.NewComponent(cfg)
	if err := sqlComponent.Initialize(); err != nil {
		log.Error("Failed to initialize SQL component", "error", err)
		os.Exit(1)
	}
	defer sqlComponent.Cleanup()

	// 初始化迁移系统
	if err := migration.InitializeMigrator(sqlComponent.DB(), cfg); err != nil {
		log.Error("Failed to initialize migration system", "error", err)
		os.Exit(1)
	}

	// 回滚迁移
	if err := migration.GlobalMigrator.Rollback(); err != nil {
		log.Error("Database migration rollback failed", "error", err)
		os.Exit(1)
	}

	log.Info("Database migration rollback completed successfully")
}

// runMigrationStatus 显示迁移状态
func runMigrationStatus(cmd *cobra.Command, args []string) {
	// 加载配置
	cfg := InitConfig(cmd, args)

	// 设置日志
	InitLogger(cfg)

	// 检查数据库是否启用
	if !cfg.Database.Enabled {
		log.Fatal("Database not enabled, please enable it in the configuration file")
	}

	// 初始化SQL组件
	sqlComponent := sqldb.NewComponent(cfg)
	if err := sqlComponent.Initialize(); err != nil {
		log.Error("Failed to initialize SQL component", "error", err)
		os.Exit(1)
	}
	defer sqlComponent.Cleanup()

	// 初始化迁移系统
	if err := migration.InitializeMigrator(sqlComponent.DB(), cfg); err != nil {
		log.Error("Failed to initialize migration system", "error", err)
		os.Exit(1)
	}

	// 获取迁移状态
	status, err := migration.GlobalMigrator.Status()
	if err != nil {
		log.Error("Failed to get migration status", "error", err)
		os.Exit(1)
	}

	// 打印迁移状态
	log.Info("Migration status:")
	for _, s := range status {
		ran := "Not executed"
		if s["ran"].(bool) {
			ran = "Executed"
		}
		log.Info("Migration",
			"version", s["version"],
			"name", s["name"],
			"status", ran,
			"batch", s["batch"],
		)
	}
}

// runMigrationReset 重置所有迁移
func runMigrationReset(cmd *cobra.Command, args []string) {
	// 加载配置
	cfg := InitConfig(cmd, args)

	// 设置日志
	InitLogger(cfg)

	// 检查数据库是否启用
	if !cfg.Database.Enabled {
		log.Fatal("Database not enabled, please enable it in the configuration file")
	}

	log.Info("Starting database migration reset process")

	// 初始化SQL组件
	sqlComponent := sqldb.NewComponent(cfg)
	if err := sqlComponent.Initialize(); err != nil {
		log.Error("Failed to initialize SQL component", "error", err)
		os.Exit(1)
	}
	defer sqlComponent.Cleanup()

	// 初始化迁移系统
	if err := migration.InitializeMigrator(sqlComponent.DB(), cfg); err != nil {
		log.Error("Failed to initialize migration system", "error", err)
		os.Exit(1)
	}

	// 重置迁移
	if err := migration.GlobalMigrator.Reset(); err != nil {
		log.Error("Database migration reset failed", "error", err)
		os.Exit(1)
	}

	log.Info("Database migration reset completed successfully")
}
