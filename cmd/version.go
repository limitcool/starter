package cmd

import (
	"fmt"

	"github.com/limitcool/starter/internal/version"
	"github.com/spf13/cobra"
)

// versionCmd 表示version子命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示应用程序的版本信息，包括版本号、构建时间和Git提交信息`,
	Run: func(cmd *cobra.Command, args []string) {
		vInfo := version.GetVersion()
		fmt.Printf("Version:      %s\n", vInfo.Version)
		fmt.Printf("GitCommit:    %s\n", vInfo.GitCommit)
		fmt.Printf("GitTreeState: %s\n", vInfo.GitTreeState)
		fmt.Printf("BuildDate:    %s\n", vInfo.BuildDate)
		fmt.Printf("GoVersion:    %s\n", vInfo.GoVersion)
		fmt.Printf("Compiler:     %s\n", vInfo.Compiler)
		fmt.Printf("Platform:     %s\n", vInfo.Platform)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
