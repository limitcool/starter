package version

import (
	"fmt"
	"runtime"
)

// 版本信息，在编译时通过ldflags设置
var (
	// Version 版本号
	Version = "v0.1.0"
	// GitCommit Git提交哈希
	GitCommit = "unknown"
	// GitTreeState Git树状态
	GitTreeState = "unknown"
	// BuildDate 构建日期
	BuildDate = "unknown"
)

// Info 包含版本信息
type Info struct {
	Version      string `json:"version"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// String 返回格式化的版本信息
func (info Info) String() string {
	return fmt.Sprintf(
		"Version: %s\nGitCommit: %s\nGitTreeState: %s\nBuildDate: %s\nGoVersion: %s\nCompiler: %s\nPlatform: %s\n",
		info.Version,
		info.GitCommit,
		info.GitTreeState,
		info.BuildDate,
		info.GoVersion,
		info.Compiler,
		info.Platform,
	)
}

// GetVersion 返回当前版本信息
func GetVersion() Info {
	return Info{
		Version:      Version,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
