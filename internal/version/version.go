package version

import (
	"fmt"
	"runtime"
)

// 版本信息，在编译时通过ldflags设置
var (
	// Version 版本号，如 v1.0.0
	Version = "dev"
	// GitCommit Git提交哈希
	GitCommit = "unknown"
	// BuildDate 构建日期
	BuildDate = "unknown"
)

// Info 包含版本信息
type Info struct {
	Version   string `json:"version"`   // 版本号
	GitCommit string `json:"gitCommit"` // Git提交哈希
	BuildDate string `json:"buildDate"` // 构建日期
	GoVersion string `json:"goVersion"` // Go版本
	Platform  string `json:"platform"`  // 平台信息
}

// String 返回格式化的版本信息
func (info Info) String() string {
	return fmt.Sprintf(
		"Version: %s\nGitCommit: %s\nBuildDate: %s\nGoVersion: %s\nPlatform: %s",
		info.Version,
		info.GitCommit,
		info.BuildDate,
		info.GoVersion,
		info.Platform,
	)
}

// Short 返回简短的版本信息
func (info Info) Short() string {
	if info.GitCommit != "unknown" && len(info.GitCommit) > 7 {
		return fmt.Sprintf("%s (%s)", info.Version, info.GitCommit[:7])
	}
	return info.Version
}

// GetVersion 返回当前版本信息
func GetVersion() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// IsDevBuild 判断是否为开发版本
func IsDevBuild() bool {
	return Version == "dev" || GitCommit == "unknown"
}

// GetVersionString 获取版本字符串
func GetVersionString() string {
	return Version
}

// GetGitCommit 获取Git提交哈希
func GetGitCommit() string {
	return GitCommit
}

// GetBuildDate 获取构建日期
func GetBuildDate() string {
	return BuildDate
}
