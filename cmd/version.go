package cmd

import (
	"fmt"

	"github.com/limitcool/starter/internal/version"
	"github.com/spf13/cobra"
)

// versionCmd 表示version子命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display application version information, including version number, build time, and Git commit information`,
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
