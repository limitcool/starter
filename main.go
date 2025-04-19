package main

import (
	_ "go.uber.org/automaxprocs"

	"github.com/limitcool/starter/cmd"
)

func main() {
	// automaxprocs 包会自动设置 GOMAXPROCS 为可用的 CPU 数量
	// 尤其在容器环境中非常有用

	// 执行根命令
	cmd.ExecuteCmd()
}
