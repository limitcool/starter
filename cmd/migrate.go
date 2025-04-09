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
	Short: "执行数据库迁移",
	Long: `执行数据库迁移，创建或更新数据库表结构。

此命令会根据定义的模型结构自动进行数据库迁移操作，适用于开发环境
和初始化生产环境。在生产环境中使用前请确保已备份数据库。`,
	Run: runMigration,
}

// migrateRollbackCmd 表示migrate:rollback子命令
var migrateRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "回滚上一批数据库迁移",
	Long:  `回滚上一批数据库迁移操作，恢复到之前的数据库状态。`,
	Run:   runMigrationRollback,
}

// migrateStatusCmd 表示migrate:status子命令
var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "显示数据库迁移状态",
	Long:  `显示所有迁移及其执行状态。`,
	Run:   runMigrationStatus,
}

// migrateResetCmd 表示migrate:reset子命令
var migrateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "重置所有数据库迁移",
	Long:  `重置所有已执行的数据库迁移，将数据库恢复到初始状态。`,
	Run:   runMigrationReset,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateRollbackCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migrateResetCmd)

	// 添加迁移特定的标志
	migrateCmd.PersistentFlags().BoolP("fresh", "f", false, "在迁移前清空数据库（危险操作）")
}

// runMigration 执行数据库迁移
func runMigration(cmd *cobra.Command, args []string) {
	// 加载配置
	cfg := InitConfig(cmd, args)

	// 设置日志
	InitLogger(cfg)

	// 检查数据库是否启用
	if !cfg.Database.Enabled {
		log.Fatal("数据库未启用，请在配置文件中启用数据库")
	}

	log.Info("开始数据库迁移流程")

	// 初始化SQL组件
	sqlComponent := sqldb.NewComponent(cfg)
	if err := sqlComponent.Initialize(); err != nil {
		log.Error("初始化SQL组件失败", "error", err)
		os.Exit(1)
	}
	defer sqlComponent.Cleanup()

	// 初始化迁移系统
	if err := migration.InitializeMigrator(sqlComponent.GetDB(), cfg); err != nil {
		log.Error("初始化迁移系统失败", "error", err)
		os.Exit(1)
	}

	// 检查是否需要重置数据库
	fresh, _ := cmd.Flags().GetBool("fresh")
	if fresh {
		log.Warn("准备清空并重建数据库...")
		if err := migration.GlobalMigrator.Reset(); err != nil {
			if err.Error() != "没有可重置的迁移" {
				log.Error("重置数据库失败", "error", err)
				os.Exit(1)
			}
		}
	}

	// 执行迁移
	if err := migration.GlobalMigrator.Migrate(); err != nil {
		log.Error("数据库迁移失败", "error", err)
		os.Exit(1)
	}

	log.Info("数据库迁移已成功完成")
}

// runMigrationRollback 回滚数据库迁移
func runMigrationRollback(cmd *cobra.Command, args []string) {
	// 加载配置
	cfg := InitConfig(cmd, args)

	// 设置日志
	InitLogger(cfg)

	// 检查数据库是否启用
	if !cfg.Database.Enabled {
		log.Fatal("数据库未启用，请在配置文件中启用数据库")
	}

	log.Info("开始数据库迁移回滚流程")

	// 初始化SQL组件
	sqlComponent := sqldb.NewComponent(cfg)
	if err := sqlComponent.Initialize(); err != nil {
		log.Error("初始化SQL组件失败", "error", err)
		os.Exit(1)
	}
	defer sqlComponent.Cleanup()

	// 初始化迁移系统
	if err := migration.InitializeMigrator(sqlComponent.GetDB(), cfg); err != nil {
		log.Error("初始化迁移系统失败", "error", err)
		os.Exit(1)
	}

	// 回滚迁移
	if err := migration.GlobalMigrator.Rollback(); err != nil {
		log.Error("数据库迁移回滚失败", "error", err)
		os.Exit(1)
	}

	log.Info("数据库迁移回滚已成功完成")
}

// runMigrationStatus 显示迁移状态
func runMigrationStatus(cmd *cobra.Command, args []string) {
	// 加载配置
	cfg := InitConfig(cmd, args)

	// 设置日志
	InitLogger(cfg)

	// 检查数据库是否启用
	if !cfg.Database.Enabled {
		log.Fatal("数据库未启用，请在配置文件中启用数据库")
	}

	// 初始化SQL组件
	sqlComponent := sqldb.NewComponent(cfg)
	if err := sqlComponent.Initialize(); err != nil {
		log.Error("初始化SQL组件失败", "error", err)
		os.Exit(1)
	}
	defer sqlComponent.Cleanup()

	// 初始化迁移系统
	if err := migration.InitializeMigrator(sqlComponent.GetDB(), cfg); err != nil {
		log.Error("初始化迁移系统失败", "error", err)
		os.Exit(1)
	}

	// 获取迁移状态
	status, err := migration.GlobalMigrator.Status()
	if err != nil {
		log.Error("获取迁移状态失败", "error", err)
		os.Exit(1)
	}

	// 打印迁移状态
	log.Info("迁移状态:")
	for _, s := range status {
		ran := "未执行"
		if s["ran"].(bool) {
			ran = "已执行"
		}
		log.Info("迁移",
			"版本", s["version"],
			"名称", s["name"],
			"状态", ran,
			"批次", s["batch"],
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
		log.Fatal("数据库未启用，请在配置文件中启用数据库")
	}

	log.Info("开始数据库迁移重置流程")

	// 初始化SQL组件
	sqlComponent := sqldb.NewComponent(cfg)
	if err := sqlComponent.Initialize(); err != nil {
		log.Error("初始化SQL组件失败", "error", err)
		os.Exit(1)
	}
	defer sqlComponent.Cleanup()

	// 初始化迁移系统
	if err := migration.InitializeMigrator(sqlComponent.GetDB(), cfg); err != nil {
		log.Error("初始化迁移系统失败", "error", err)
		os.Exit(1)
	}

	// 重置迁移
	if err := migration.GlobalMigrator.Reset(); err != nil {
		log.Error("数据库迁移重置失败", "error", err)
		os.Exit(1)
	}

	log.Info("数据库迁移重置已成功完成")
}
