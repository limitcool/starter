package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout  = 30 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

// Component MongoDB组件实现
type Component struct {
	Config  *configs.Config
	client  *mongo.Client
	enabled bool
}

// NewComponent 创建MongoDB组件
func NewComponent(cfg *configs.Config) *Component {
	return &Component{
		Config:  cfg,
		enabled: cfg.Mongo.Enabled && cfg.Mongo.URI != "" && cfg.Mongo.DB != "",
	}
}

// Name 返回组件名称
func (m *Component) Name() string {
	return "MongoDB"
}

// Initialize 初始化MongoDB连接
func (m *Component) Initialize() error {
	if !m.enabled {
		logger.Info("MongoDB component disabled")
		return nil
	}

	logger.Info("Initializing MongoDB component")

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.Config.Mongo.URI).
		SetAuth(options.Credential{
			Username: m.Config.Mongo.User,
			Password: m.Config.Mongo.Password,
		}).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize))

	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 验证连接
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	m.client = client

	// 不再设置全局实例
	// setupInstance(m) // 已移除，使用依赖注入代替

	// 不再设置兼容性全局变量
	// Mongo = client // 已移除，使用依赖注入代替

	logger.Info("MongoDB component initialized successfully")
	return nil
}

// Cleanup 清理MongoDB资源
func (m *Component) Cleanup() {
	if m.client != nil {
		logger.Info("Closing MongoDB connection")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := m.client.Disconnect(ctx); err != nil {
			logger.Error("Error disconnecting from MongoDB", "error", err)
		}
	}
}

// IsEnabled 检查组件是否启用
func (m *Component) IsEnabled() bool {
	return m.enabled
}

// GetClient 获取MongoDB客户端
func (m *Component) GetClient() *mongo.Client {
	return m.client
}

// GetDB 获取默认数据库
func (m *Component) GetDB() *mongo.Database {
	if m.client == nil {
		return nil
	}
	return m.client.Database(m.Config.Mongo.DB)
}

// GetCollection 获取集合
func (m *Component) GetCollection(name string) *mongo.Collection {
	if m.client == nil {
		return nil
	}
	return m.client.Database(m.Config.Mongo.DB).Collection(name)
}
