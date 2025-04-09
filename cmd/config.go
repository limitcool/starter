package cmd

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/pkg/env"
	"github.com/limitcool/starter/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 全局配置实例
var globalConfig *configs.Config

// InitConfig 加载配置文件
func InitConfig(cmd *cobra.Command, args []string) *configs.Config {
	if globalConfig != nil {
		return globalConfig
	}

	// 先设置基本日志格式，确保在配置读取前就使用统一格式
	log.SetDefault(log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.InfoLevel,
		Prefix:          "🌏 starter",
		TimeFormat:      time.RFC3339,
		ReportTimestamp: true,
		Formatter:       log.TextFormatter,
	}))

	// 检查是否通过flag指定了配置文件
	configFile, _ := cmd.Flags().GetString("config")

	// 如果未指定配置文件，则使用环境名称
	if configFile == "" {
		envName := env.Get().String()
		configName := envName // 使用环境名称作为配置文件名: dev.yaml, test.yaml, prod.yaml

		viper.SetConfigName(configName)
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.SetConfigType("yaml")
	} else {
		// 使用指定的配置文件
		viper.SetConfigFile(configFile)

		// 设置配置类型（基于文件扩展名）
		ext := filepath.Ext(configFile)
		if ext != "" {
			viper.SetConfigType(ext[1:]) // 去除前导点号
		}
	}

	// 读取环境配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config file", "error", err)
	}

	// 输出使用的配置文件
	log.Info("Using config file", "path", viper.ConfigFileUsed())

	// 解析配置到结构体
	cfg := &configs.Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatal("Config unmarshal failed", "error", err)
	}

	// 初始化服务管理器
	services.Init(cfg, nil)
	globalConfig = cfg

	return cfg
}

// InitLogger 配置全局日志
func InitLogger(cfg *configs.Config) {
	log.SetPrefix("🌏 starter ")

	// 获取环境
	currentEnv := env.Get()

	// 根据环境设置Gin模式
	if currentEnv == env.Dev {
		// 在开发环境中，我们可以保留Gin的调试输出
		gin.SetMode(gin.DebugMode)

		// 但仍然将它重定向到我们的日志系统
		logger.SetupGinLogger()
	} else {
		// 在非开发环境中，完全禁用Gin的调试输出
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
	}

	// 使用配置更新日志设置
	logger.Setup(cfg.Log)

	// 记录环境模式
	if env.IsProd() {
		log.Info("Running in production mode")
	} else if env.IsTest() {
		log.Info("Running in test mode")
	} else {
		log.Info("Running in debug mode")
	}
}
