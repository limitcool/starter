package main

import (
	"github.com/limitcool/starter/cmd"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	// 静默设置 GOMAXPROCS，不输出日志
	_, _ = maxprocs.Set(maxprocs.Logger(func(string, ...any) {}))

	// 执行根命令
	cmd.ExecuteCmd()
}
