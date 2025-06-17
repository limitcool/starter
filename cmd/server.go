package cmd

import (
	"github.com/limitcool/starter/internal/app"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/version"
	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
)

// serverCmd 表示server子命令
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP server",
	Long: `Start HTTP server and provide Web API services.

The server will load the configuration file and initialize all necessary components, including database connections, logging systems, etc.
The server gracefully handles shutdown signals, ensuring all requests are processed and resources are safely closed.`,
	Run: runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// 添加服务器特定的标志
	serverCmd.Flags().IntP("port", "p", 0, "HTTP server port number, overrides the setting in the configuration file")
}

// runServer 运行HTTP服务器
func runServer(cmd *cobra.Command, args []string) {
	// 静默设置 GOMAXPROCS，避免日志输出
	_, _ = maxprocs.Set(maxprocs.Logger(func(string, ...any) {}))

	// 加载配置
	cfg := InitConfig(cmd, args)

	// 设置日志
	InitLogger(cfg)

	// 显示版本信息
	vInfo := version.GetVersion()
	gitCommitShort := vInfo.GitCommit
	if len(gitCommitShort) > 8 {
		gitCommitShort = gitCommitShort[:8]
	}
	logger.Info("Application starting",
		"version", vInfo.Version,
		"gitCommit", gitCommitShort,
		"buildDate", vInfo.BuildDate,
		"goVersion", vInfo.GoVersion,
		"platform", vInfo.Platform)

	// 检查是否从命令行指定了端口
	port, _ := cmd.Flags().GetInt("port")
	if port > 0 {
		cfg.App.Port = port
	}

	logger.Info("Application starting with manual dependency injection")

	// 创建应用实例
	application, err := app.New(cfg)
	if err != nil {
		logger.Error("Failed to create application", "error", err)
		return
	}

	// 运行应用
	if err := application.Run(); err != nil {
		logger.Error("Application run failed", "error", err)
	}
}
