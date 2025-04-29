package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// 运行所有测试
func main() {
	// 获取项目根目录
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	// 切换到项目根目录
	rootDir = filepath.Dir(rootDir)
	if err := os.Chdir(rootDir); err != nil {
		fmt.Printf("切换到项目根目录失败: %v\n", err)
		os.Exit(1)
	}

	// 运行单元测试
	fmt.Println("运行单元测试...")
	cmd := exec.Command("go", "test", "./internal/pkg/casbin", "./internal/middleware", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("单元测试失败: %v\n", err)
		os.Exit(1)
	}

	// 运行集成测试
	fmt.Println("\n运行集成测试...")
	cmd = exec.Command("go", "test", "./test/integration", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("集成测试失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n所有测试通过!")
}
