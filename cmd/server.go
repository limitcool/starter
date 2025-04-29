package cmd

import (
	"context"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/api"
	"github.com/limitcool/starter/internal/datastore/redisdb"
	"github.com/limitcool/starter/internal/datastore/sqldb"
	"github.com/limitcool/starter/internal/filestore"
	"github.com/limitcool/starter/internal/handler"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/router"
	"github.com/limitcool/starter/internal/server"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
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

// ConfigParams 配置参数
type ConfigParams struct {
	Cmd  *cobra.Command
	Args []string
}

// LoadConfig 加载配置
func LoadConfig(params *ConfigParams) (*configs.Config, error) {
	// 加载配置
	cfg := InitConfig(params.Cmd, params.Args)

	// 设置日志
	InitLogger(cfg)

	// 检查是否从命令行指定了端口
	port, _ := params.Cmd.Flags().GetInt("port")
	if port > 0 {
		cfg.App.Port = port
	}

	return cfg, nil
}

// runServer 运行HTTP服务器
func runServer(cmd *cobra.Command, args []string) {
	// 创建fx应用程序
	app := fx.New(
		// 提供命令行参数
		fx.Supply(cmd, args),

		// 提供配置加载函数
		fx.Provide(
			func(cmd *cobra.Command, args []string) *ConfigParams {
				return &ConfigParams{
					Cmd:  cmd,
					Args: args,
				}
			},
			LoadConfig,
		),

		// 添加所有模块 - 在lite版本中，只添加必要的模块
		sqldb.Module,
		redisdb.Module,
		filestore.Module,
		// 使用handler模块替代repository、service和controller模块
		handler.Module,
		api.Module,
		router.Module,
		server.Module,

		// 注册生命周期钩子
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info("Application started with fx framework")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("Application stopped")
					return nil
				},
			})
		}),
	)

	// 启动应用
	app.Run()
}
