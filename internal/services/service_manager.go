package services

import (
	"sync"

	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"gorm.io/gorm"
)

// ServiceManager 服务管理器，管理所有服务的单例实例及核心资源
type ServiceManager struct {
	// 核心资源
	db     *gorm.DB
	config *configs.Config

	// 各种服务的单例实例
	userService         *SysUserService
	roleService         *RoleService
	casbinService       *CasbinService
	operationLogService *OperationLogService
	normalUserService   *UserService

	// 其他服务可以按需添加
}

var (
	serviceInstance *ServiceManager
	serviceOnce     sync.Once
)

// Instance 获取ServiceManager的单例实例
func Instance() *ServiceManager {
	// 初始化实例但不注入资源，资源将通过Init方法注入
	serviceOnce.Do(func() {
		serviceInstance = &ServiceManager{}
	})
	return serviceInstance
}

// Init 初始化服务管理器并注入资源
// 这个方法应该在应用启动时调用，在所有组件初始化完成后
func Init(config *configs.Config, db *gorm.DB) {
	mgr := Instance()

	// 注入核心资源
	mgr.config = config
	mgr.db = db
}

// GetDB 获取数据库连接
func (sm *ServiceManager) GetDB() *gorm.DB {
	// 优先使用注入的DB，如果未注入则尝试从sqldb包获取
	if sm.db != nil {
		return sm.db
	}
	return sqldb.GetDB()
}

// GetConfig 获取应用配置
func (sm *ServiceManager) GetConfig() *configs.Config {
	if sm.config == nil {
		panic("配置未初始化")
	}
	return sm.config
}

// GetUserService 获取UserService实例
func (sm *ServiceManager) GetUserService() *SysUserService {
	if sm.userService == nil {
		sm.userService = NewSysUserService()
	}
	return sm.userService
}

// GetRoleService 获取RoleService实例
func (sm *ServiceManager) GetRoleService() *RoleService {
	if sm.roleService == nil {
		sm.roleService = NewRoleService()
	}
	return sm.roleService
}

// GetCasbinService 获取CasbinService实例
func (sm *ServiceManager) GetCasbinService() *CasbinService {
	if sm.casbinService == nil {
		sm.casbinService = NewCasbinService()
	}
	return sm.casbinService
}

// GetOperationLogService 获取OperationLogService实例
func (sm *ServiceManager) GetOperationLogService() *OperationLogService {
	if sm.operationLogService == nil {
		sm.operationLogService = NewOperationLogService()
	}
	return sm.operationLogService
}

// GetNormalUserService 获取NormalUserService实例
func (sm *ServiceManager) GetNormalUserService() *UserService {
	if sm.normalUserService == nil {
		sm.normalUserService = NewUserService()
	}
	return sm.normalUserService
}
