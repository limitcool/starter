package api

import (
	"github.com/limitcool/starter/internal/controller"
	"go.uber.org/fx"
)

// Module API模块
var Module = fx.Options(
	// 包含控制器模块
	controller.Module,
)
