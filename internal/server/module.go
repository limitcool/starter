package server

import (
	"github.com/limitcool/starter/internal/server/grpc"
	"github.com/limitcool/starter/internal/server/http"
	"go.uber.org/fx"
)

// Module 服务器模块
var Module = fx.Options(
	// 提供HTTP服务器
	fx.Provide(http.NewHTTPServer),
	
	// 提供gRPC服务器
	fx.Provide(grpc.NewGRPCServer),
)
