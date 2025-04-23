package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 定义需要替换的日志函数
var loggerFunctions = []string{
	"logger.Info",
	"logger.Error",
	"logger.Warn",
	"logger.Debug",
	"logger.Fatal",
}

// 定义对应的带上下文的函数
var contextLoggerFunctions = []string{
	"logger.InfoContext",
	"logger.ErrorContext",
	"logger.WarnContext",
	"logger.DebugContext",
	"logger.FatalContext",
}

// 定义需要扫描的文件扩展名
var fileExtensions = []string{".go"}

// 定义需要排除的目录
var excludeDirs = []string{
	"vendor",
	"node_modules",
	".git",
	"scripts",
	"internal/pkg/logger", // 排除日志包本身
}

// 检查目录是否应该被排除
func shouldExcludeDir(dir string) bool {
	for _, excludeDir := range excludeDirs {
		if strings.Contains(dir, excludeDir) {
			return true
		}
	}
	return false
}

// 检查文件是否应该被处理
func shouldProcessFile(file string) bool {
	ext := filepath.Ext(file)
	for _, fileExt := range fileExtensions {
		if ext == fileExt {
			return true
		}
	}
	return false
}

// 替换日志函数调用
func replaceLoggerCalls(content string) string {
	// 创建正则表达式来匹配日志函数调用
	for i, logFunc := range loggerFunctions {
		// 匹配不带上下文的日志调用，但不匹配已经带上下文的调用
		// 例如：logger.Info("message", ...) 但不匹配 logger.InfoContext(ctx, "message", ...)
		pattern := regexp.MustCompile(fmt.Sprintf(`%s\(([^)]*)\)`, regexp.QuoteMeta(logFunc)))

		// 替换为带上下文的版本
		content = pattern.ReplaceAllStringFunc(content, func(match string) string {
			// 如果已经是带上下文的调用，不做替换
			if strings.Contains(match, "Context") {
				return match
			}

			// 提取参数
			paramsStart := strings.Index(match, "(")
			paramsEnd := strings.LastIndex(match, ")")
			if paramsStart == -1 || paramsEnd == -1 {
				return match
			}

			params := match[paramsStart+1 : paramsEnd]

			// 构建新的调用
			return fmt.Sprintf("%s(ctx, %s)", contextLoggerFunctions[i], params)
		})
	}

	return content
}

// 处理单个文件
func processFile(filePath string) error {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 检查文件是否包含 context 导入
	hasContext := strings.Contains(string(content), "context")

	// 检查文件是否包含日志函数调用
	hasLoggerCall := false
	for _, logFunc := range loggerFunctions {
		if strings.Contains(string(content), logFunc) {
			hasLoggerCall = true
			break
		}
	}

	// 如果文件包含日志调用但没有 context 导入，可能需要手动处理
	if hasLoggerCall && !hasContext {
		fmt.Printf("警告: 文件 %s 包含日志调用但可能没有 context 导入，需要手动检查\n", filePath)
		return nil
	}

	// 替换日志函数调用
	newContent := replaceLoggerCalls(string(content))

	// 如果内容有变化，写回文件
	if newContent != string(content) {
		fmt.Printf("更新文件: %s\n", filePath)
		return os.WriteFile(filePath, []byte(newContent), 0644)
	}

	return nil
}

// 遍历目录处理文件
func walkDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过排除的目录
		if info.IsDir() && shouldExcludeDir(path) {
			return filepath.SkipDir
		}

		// 处理符合条件的文件
		if !info.IsDir() && shouldProcessFile(path) {
			if err := processFile(path); err != nil {
				fmt.Printf("处理文件 %s 时出错: %v\n", path, err)
			}
		}

		return nil
	})
}

func main() {
	// 获取要处理的目录
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	fmt.Printf("开始处理目录: %s\n", dir)

	// 遍历目录处理文件
	if err := walkDir(dir); err != nil {
		fmt.Printf("处理目录时出错: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("处理完成")
}
