package grpc

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// ServerParams gRPC服务器参数
type ServerParams struct {
	fx.In

	LC     fx.Lifecycle
	Config *configs.Config
}

// 全局监听器变量
var grpcListener net.Listener

// NewGRPCServer 创建gRPC服务器
func NewGRPCServer(params ServerParams) *grpc.Server {
	// 如果gRPC服务未启用，返回nil
	if !params.Config.GRPC.Enabled {
		logger.Info("gRPC server is disabled")
		return nil
	}

	// 创建gRPC服务器选项
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// 添加恢复中间件
			grpc_recovery.UnaryServerInterceptor(),
			// 添加日志中间件
			LoggingInterceptor(),
			// 添加错误处理中间件
			ErrorHandlerInterceptor(),
			// 添加请求上下文中间件
			ContextInterceptor(),
		)),
	}

	// 创建gRPC服务器
	server := grpc.NewServer(opts...)

	// 创建监听器
	var err error
	grpcListener, err = net.Listen("tcp", fmt.Sprintf(":%d", params.Config.GRPC.Port))
	if err != nil {
		logger.Error("Failed to listen", "error", err)
		return nil
	}

	// 注册健康检查服务
	if params.Config.GRPC.HealthCheck {
		healthServer := health.NewServer()
		grpc_health_v1.RegisterHealthServer(server, healthServer)
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	}

	// 注册反射服务
	if params.Config.GRPC.Reflection {
		reflection.Register(server)
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("==================================================")
			logger.Info("gRPC服务器已启动", 
				"address", fmt.Sprintf("localhost:%d", params.Config.GRPC.Port),
				"reflection", params.Config.GRPC.Reflection,
				"health_check", params.Config.GRPC.HealthCheck)
			
			// 打印gRPC服务信息
			printGRPCServices(server)
			
			go func() {
				if err := server.Serve(grpcListener); err != nil {
					logger.Error("gRPC server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping gRPC server")
			server.GracefulStop()
			logger.Info("gRPC server stopped")
			return nil
		},
	})

	return server
}

// printGRPCServices 打印所有gRPC服务信息
func printGRPCServices(server *grpc.Server) {
	if server == nil {
		return
	}

	// 获取所有服务信息
	serviceInfo := server.GetServiceInfo()
	if len(serviceInfo) == 0 {
		return
	}

	// 打印服务信息
	logger.Info("gRPC服务列表:")
	
	for name, info := range serviceInfo {
		// 跳过reflection服务
		if strings.HasPrefix(name, "grpc.reflection") {
			continue
		}
		
		logger.Info("服务名称", "name", name)
		for _, method := range info.Methods {
			logger.Info("  - 方法", "name", method.Name)
		}
	}
	
	logger.Info("==================================================")
}

// LoggingInterceptor 日志中间件
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		
		// 记录请求
		logger.InfoContext(ctx, "gRPC request started", 
			"method", info.FullMethod,
			"request", req,
		)
		
		// 处理请求
		resp, err := handler(ctx, req)
		
		// 记录响应
		logger.InfoContext(ctx, "gRPC request completed",
			"method", info.FullMethod,
			"duration", time.Since(start),
			"error", err,
		)
		
		return resp, err
	}
}

// ErrorHandlerInterceptor 错误处理中间件
func ErrorHandlerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			// 记录错误
			logger.ErrorContext(ctx, "gRPC error",
				"method", info.FullMethod,
				"error", err,
			)
		}
		return resp, err
	}
}

// ContextInterceptor 请求上下文中间件
func ContextInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 生成请求ID
		requestID := fmt.Sprintf("grpc-req-%d", time.Now().UnixNano())
		
		// 添加请求ID到上下文
		ctx = context.WithValue(ctx, "request_id", requestID)
		
		// 处理请求
		return handler(ctx, req)
	}
}
