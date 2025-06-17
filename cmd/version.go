package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/limitcool/starter/internal/version"
	"github.com/spf13/cobra"
)

var (
	// version命令的标志
	outputFormat string
)

// versionCmd 表示version子命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display application version information, including version number, build time, and Git commit information`,
	Run: func(cmd *cobra.Command, args []string) {
		vInfo := version.GetVersion()

		switch outputFormat {
		case "json":
			// JSON格式输出
			jsonData, err := json.MarshalIndent(vInfo, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting JSON: %v\n", err)
				return
			}
			fmt.Println(string(jsonData))
		case "short":
			// 简短格式输出
			fmt.Println(vInfo.Short())
		default:
			// 默认详细格式输出
			fmt.Println(vInfo.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// 添加输出格式标志
	versionCmd.Flags().StringVarP(&outputFormat, "output", "o", "default", "Output format (default|json|short)")
}
