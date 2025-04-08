package cmd

import (
	"os"

	"github.com/charmbracelet/log"
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

func init() {
	rootCmd.AddCommand(migrateCmd)

	// 添加迁移特定的标志
	migrateCmd.Flags().BoolP("reset", "r", false, "重置数据库（危险操作）")
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

	// 执行迁移
	if err := sqlComponent.Migrate(); err != nil {
		log.Error("数据库迁移失败", "error", err)
		os.Exit(1)
	}

	log.Info("数据库迁移已成功完成")
}
