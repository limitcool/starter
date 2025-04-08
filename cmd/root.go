package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "starter",
	Short: "Starter是一个Go Web应用程序框架",
	Long: `Starter是一个Go Web应用程序框架，提供了完整的Web应用骨架，
包括配置管理、数据库集成、日志系统等。

使用Starter可以快速构建高质量的Web应用，专注于业务逻辑而非基础设施。
详情请访问: https://github.com/limitcool/starter`,
}

// ExecuteCmd 执行rootCmd
func ExecuteCmd() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 在这里可以设置全局标志（flags）
	rootCmd.PersistentFlags().StringP("config", "c", "", "配置文件路径")

	// 设置版本信息
	rootCmd.Version = "1.0.0"
}
