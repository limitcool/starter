package services

import (
	"sync"

	"github.com/limitcool/starter/internal/repository"
	"github.com/limitcool/starter/internal/storage/database"
)

// Factory 服务工厂
// 用于创建和管理服务实例
type Factory struct {
	db database.DB
	mu sync.RWMutex

	// 仓库实例
	menuRepo       *repository.MenuRepo
	roleRepo       *repository.RoleRepo
	userRepo       *repository.GormUserRepository
	sysUserRepo    *repository.SysUserRepo
	permissionRepo *repository.PermissionRepo
	// fileRepo         repository.FileRepository // 暂时注释，因为没有找到具体实现类
	operationLogRepo *repository.OperationLogRepo

	// 服务实例缓存
	sysUserService      *SysUserService
	roleService         *RoleService
	menuService         *MenuService
	permissionService   *PermissionService
	operationLogService *OperationLogService
	fileService         *FileService
	casbinService       *CasbinService
	systemService       *SystemService
}

// NewFactory 创建服务工厂
func NewFactory(db database.DB) *Factory {
	// 获取GORM DB实例
	gormDB := db.GetDB()

	// 创建仓库实例
	menuRepo := repository.NewMenuRepo(gormDB)
	roleRepo := repository.NewRoleRepo(gormDB)
	sysUserRepo := repository.NewSysUserRepo(gormDB)
	userRepo := repository.NewUserRepository(gormDB)

	return &Factory{
		db:          db,
		menuRepo:    menuRepo,
		roleRepo:    roleRepo,
		sysUserRepo: sysUserRepo,
		userRepo:    userRepo,
	}
}

// SysUser 获取系统用户服务
func (f *Factory) SysUser() *SysUserService {
	f.mu.RLock()
	if f.sysUserService != nil {
		defer f.mu.RUnlock()
		return f.sysUserService
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.sysUserService == nil {
		// 创建依赖服务
		roleService := f.Role()

		// 创建用户服务
		f.sysUserService = NewSysUserService(f.sysUserRepo, f.userRepo, roleService)
	}
	return f.sysUserService
}

// Role 获取角色服务
func (f *Factory) Role() *RoleService {
	f.mu.RLock()
	if f.roleService != nil {
		defer f.mu.RUnlock()
		return f.roleService
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.roleService == nil {
		// 创建依赖服务
		casbinService := f.Casbin()

		// 创建角色服务
		f.roleService = NewRoleService(f.roleRepo, casbinService)
	}
	return f.roleService
}

// Menu 获取菜单服务
func (f *Factory) Menu() *MenuService {
	f.mu.RLock()
	if f.menuService != nil {
		defer f.mu.RUnlock()
		return f.menuService
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.menuService == nil {
		// 创建依赖服务
		casbinService := f.Casbin()

		// 创建菜单服务
		f.menuService = NewMenuService(f.menuRepo, casbinService)
	}
	return f.menuService
}

// Permission 获取权限服务
func (f *Factory) Permission() *PermissionService {
	f.mu.RLock()
	if f.permissionService != nil {
		defer f.mu.RUnlock()
		return f.permissionService
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.permissionService == nil {
		// 创建权限服务
		f.permissionService = NewPermissionService(f.permissionRepo)
	}
	return f.permissionService
}

// OperationLog 获取操作日志服务
func (f *Factory) OperationLog() *OperationLogService {
	f.mu.RLock()
	if f.operationLogService != nil {
		defer f.mu.RUnlock()
		return f.operationLogService
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.operationLogService == nil {
		// 创建操作日志服务
		f.operationLogService = NewOperationLogService(f.operationLogRepo)
	}
	return f.operationLogService
}

// File 获取文件服务
func (f *Factory) File() *FileService {
	f.mu.RLock()
	if f.fileService != nil {
		defer f.mu.RUnlock()
		return f.fileService
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fileService == nil {
		// 创建文件服务
		// 注意: 需要实现NewFileService(db database.DB)
		// f.fileService = NewFileService(f.db)
	}
	return f.fileService
}

// Casbin 获取Casbin服务
func (f *Factory) Casbin() *CasbinService {
	f.mu.RLock()
	if f.casbinService != nil {
		defer f.mu.RUnlock()
		return f.casbinService
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.casbinService == nil {
		// 创建Casbin服务
		f.casbinService = NewCasbinService(f.db)
	}
	return f.casbinService
}

// System 获取系统服务
func (f *Factory) System() *SystemService {
	f.mu.RLock()
	if f.systemService != nil {
		defer f.mu.RUnlock()
		return f.systemService
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.systemService == nil {
		f.systemService = NewSystemService(f.db)
	}
	return f.systemService
}

// 全局服务工厂实例
var (
	globalFactory *Factory
	factoryOnce   sync.Once
)

// SetGlobalFactory 设置全局服务工厂
func SetGlobalFactory(db database.DB) {
	factoryOnce.Do(func() {
		globalFactory = NewFactory(db)
	})
}

// GetFactory 获取全局服务工厂
func GetFactory() *Factory {
	return globalFactory
}
