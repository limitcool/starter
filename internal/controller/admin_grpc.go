package controller

import (
	"context"
	"time"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	pb "github.com/limitcool/starter/internal/proto/gen/v1"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/internal/version"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// AdminGRPCController gRPC管理系统控制器
type AdminGRPCController struct {
	pb.UnimplementedSystemServiceServer
	config       *configs.Config
	adminService *services.AdminService
}

// NewAdminGRPCController 创建gRPC管理系统控制器
func NewAdminGRPCController(config *configs.Config, adminService *services.AdminService) *AdminGRPCController {
	return &AdminGRPCController{
		config:       config,
		adminService: adminService,
	}
}

// GetSystemInfo 获取系统信息
func (c *AdminGRPCController) GetSystemInfo(ctx context.Context, req *pb.SystemInfoRequest) (*pb.SystemInfoResponse, error) {
	// 记录请求
	logger.InfoContext(ctx, "GetSystemInfo called", "request_id", req.RequestId)

	// 获取版本信息
	verInfo := version.GetVersion()

	// 构建响应
	resp := &pb.SystemInfoResponse{
		AppName:    c.config.App.Name,
		Version:    verInfo.Version,
		Mode:       c.config.App.Mode,
		ServerTime: time.Now().Unix(),
	}

	return resp, nil
}

// AdminGRPCControllerParams gRPC管理系统控制器参数
type AdminGRPCControllerParams struct {
	fx.In

	LC           fx.Lifecycle
	Config       *configs.Config
	AdminService *services.AdminService
	GRPCServer   *grpc.Server `optional:"true"`
}

// RegisterAdminGRPCController 注册gRPC管理系统控制器
func RegisterAdminGRPCController(params AdminGRPCControllerParams) {
	// 如果gRPC服务未启用或服务器为nil，直接返回
	if !params.Config.GRPC.Enabled || params.GRPCServer == nil {
		return
	}

	// 创建控制器
	controller := NewAdminGRPCController(params.Config, params.AdminService)

	// 注册服务
	pb.RegisterSystemServiceServer(params.GRPCServer, controller)

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("gRPC管理系统服务已注册", "service", "AdminService")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("gRPC管理系统服务已停止", "service", "AdminService")
			return nil
		},
	})
}
