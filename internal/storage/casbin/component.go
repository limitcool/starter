package casbin

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"gorm.io/gorm"
)

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

// Component Casbin组件
type Component struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
	config   *configs.Config
}

// NewComponent 创建Casbin组件
func NewComponent(cfg *configs.Config) *Component {
	return &Component{
		config: cfg,
	}
}

// Initialize 初始化Casbin组件
func (c *Component) Initialize() error {
	// 如果权限系统未启用，直接返回
	if c.config != nil && !c.config.Permission.Enabled {
		return nil
	}

	// 获取数据库连接
	c.db = services.Instance().GetDB()
	if c.db == nil {
		// 尝试从SQL组件获取数据库连接
		c.db = sqldb.GetDB()
		if c.db == nil {
			return fmt.Errorf("数据库未初始化")
		}
	}

	// 初始化Casbin
	e, initErr := c.initEnforcer()
	if initErr != nil {
		return initErr
	}

	c.enforcer = e
	enforcer = e // 设置全局enforcer

	return nil
}

// 初始化Casbin Enforcer
func (c *Component) initEnforcer() (*casbin.Enforcer, error) {
	var err error
	var e *casbin.Enforcer

	once.Do(func() {
		// 使用gorm适配器
		adapter, adapterErr := gormadapter.NewAdapterByDB(c.db)
		if adapterErr != nil {
			err = adapterErr
			return
		}

		// 获取模型文件路径，如果配置中没有指定则使用默认值
		modelPath := "configs/rbac_model.conf"
		if c.config != nil && c.config.Permission.ModelPath != "" {
			modelPath = c.config.Permission.ModelPath
		}

		// 创建enforcer
		e, err = casbin.NewEnforcer(modelPath, adapter)
		if err != nil {
			return
		}

		// 加载策略
		if loadErr := e.LoadPolicy(); loadErr != nil {
			err = loadErr
			return
		}

		// 启用自动保存
		e.EnableAutoSave(true)
	})

	return e, err
}

// GetEnforcer 获取Casbin实例
func (c *Component) GetEnforcer() *casbin.Enforcer {
	return c.enforcer
}

// 获取全局Enforcer实例
func GetEnforcer() *casbin.Enforcer {
	// 如果权限系统未启用，直接返回nil
	svcMgr := services.Instance()
	if svcMgr.GetConfig() != nil && !svcMgr.GetConfig().Permission.Enabled {
		return nil
	}
	return enforcer
}

// Cleanup 清理资源
func (c *Component) Cleanup() error {
	return nil
}

// Migrate 执行Casbin相关迁移
func (c *Component) Migrate() error {
	// Casbin适配器会自动创建必要的表，这里不需要额外操作
	log.Info("Casbin表迁移完成")
	return nil
}

// CheckPermission 检查权限
func (c *Component) CheckPermission(userID string, obj string, act string) (bool, error) {
	return c.enforcer.Enforce(userID, obj, act)
}

// AddRoleForUser 为用户添加角色
func (c *Component) AddRoleForUser(userID string, role string) (bool, error) {
	return c.enforcer.AddGroupingPolicy(userID, role)
}

// DeleteRoleForUser 删除用户角色
func (c *Component) DeleteRoleForUser(userID string, role string) (bool, error) {
	return c.enforcer.RemoveGroupingPolicy(userID, role)
}

// AddPermissionForRole 为角色添加权限
func (c *Component) AddPermissionForRole(role string, obj string, act string) (bool, error) {
	return c.enforcer.AddPolicy(role, obj, act)
}

// DeletePermissionForRole 删除角色权限
func (c *Component) DeletePermissionForRole(role string, obj string, act string) (bool, error) {
	return c.enforcer.RemovePolicy(role, obj, act)
}

// GetRolesForUser 获取用户角色列表
func (c *Component) GetRolesForUser(userID string) ([]string, error) {
	return c.enforcer.GetRolesForUser(userID)
}

// GetPermissionsForUser 获取用户所有权限
func (c *Component) GetPermissionsForUser(userID string) [][]string {
	permissions, _ := c.enforcer.GetImplicitPermissionsForUser(userID)
	return permissions
}

// GetPermissionsForRole 获取角色所有权限
func (c *Component) GetPermissionsForRole(role string) [][]string {
	permissions, _ := c.enforcer.GetFilteredPolicy(0, role)
	return permissions
}

// DeleteRole 删除角色
func (c *Component) DeleteRole(role string) (bool, error) {
	// 删除角色所有权限
	_, err := c.enforcer.RemoveFilteredPolicy(0, role)
	if err != nil {
		return false, err
	}

	// 删除所有用户与该角色的关联
	_, err = c.enforcer.RemoveFilteredGroupingPolicy(1, role)
	if err != nil {
		return false, err
	}

	return true, nil
}
