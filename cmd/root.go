package cmd

import (
	"fmt"
	"os"

	"github.com/limitcool/starter/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "starter",
	Short: "Starter is a Go Web application framework",
	Long: `Starter is a Go Web application framework that provides a complete Web application skeleton,
including configuration management, database integration, logging system, etc.

With Starter, you can quickly build high-quality Web applications, focusing on business logic rather than infrastructure.
For details, please visit: https://github.com/limitcool/starter`,
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
	rootCmd.PersistentFlags().StringP("config", "c", "", "Configuration file path")

	// 设置版本信息
	verInfo := version.GetVersion()
	rootCmd.Version = verInfo.Version

	// 自定义版本模板
	rootCmd.SetVersionTemplate(fmt.Sprintf("{{.Name}} %s\n\n%s", verInfo.Version, verInfo.String()))
}
