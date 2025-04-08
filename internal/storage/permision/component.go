package permission

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/storage/casbin"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"gorm.io/gorm"
)

// Component 权限系统组件
type Component struct {
	db         *gorm.DB
	config     *configs.Config
	casbinComp *casbin.Component
}

// NewComponent 创建权限系统组件
func NewComponent(cfg *configs.Config) *Component {
	return &Component{
		config: cfg,
	}
}

// Initialize 初始化权限系统组件
func (c *Component) Initialize() error {
	// 优先使用sqldb.GetDB()获取数据库连接
	c.db = sqldb.GetDB()
	if c.db == nil {
		// 尝试从SQL组件获取数据库连接
		sqlComp := sqldb.NewComponent(c.config)
		if err := sqlComp.Initialize(); err != nil {
			return err
		}
		c.db = sqldb.GetDB()

		// 再次检查数据库连接
		if c.db == nil {
			return fmt.Errorf("数据库未初始化")
		}
	}

	// 初始化Casbin组件
	c.casbinComp = casbin.NewComponent(c.config)
	if err := c.casbinComp.Initialize(); err != nil {
		return err
	}

	return nil
}

// Cleanup 清理资源
func (c *Component) Cleanup() error {
	if c.casbinComp != nil {
		return c.casbinComp.Cleanup()
	}
	return nil
}

// Migrate 执行权限系统相关迁移
func (c *Component) Migrate() error {
	// 迁移权限相关表
	if err := c.db.AutoMigrate(
		&model.Menu{},
		&model.Role{},
		&model.RoleMenu{},
		&model.UserRole{},
		&model.Permission{},
		&model.RolePermission{},
	); err != nil {
		return err
	}

	// 执行Casbin表迁移
	if err := c.casbinComp.Migrate(); err != nil {
		return err
	}

	// 初始化基础数据
	if err := c.initBaseData(); err != nil {
		return err
	}

	log.Info("权限系统表迁移完成")
	return nil
}

// 初始化基础数据
func (c *Component) initBaseData() error {
	// 初始化管理员角色
	if err := c.initAdminRole(); err != nil {
		return err
	}

	// 初始化基础菜单
	if err := c.initBasicMenus(); err != nil {
		return err
	}

	return nil
}

// 初始化管理员角色
func (c *Component) initAdminRole() error {
	// 检查是否已有admin角色
	var count int64
	if err := c.db.Model(&model.Role{}).Where("code = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	// 已存在则不重复创建
	if count > 0 {
		log.Info("管理员角色已存在，跳过创建")
		return nil
	}

	// 创建超级管理员角色
	adminRole := model.Role{
		Name:        "超级管理员",
		Code:        "admin",
		Enabled:     true,
		Sort:        0,
		Description: "系统超级管理员",
	}

	if err := c.db.Create(&adminRole).Error; err != nil {
		return err
	}

	// 为超级管理员设置全部权限(使用通配符)
	if _, err := c.casbinComp.AddPermissionForRole("admin", "*", "*"); err != nil {
		return err
	}

	log.Info("管理员角色创建成功")
	return nil
}

// 初始化基础菜单
func (c *Component) initBasicMenus() error {
	// 检查是否已有菜单
	var count int64
	if err := c.db.Model(&model.Menu{}).Count(&count).Error; err != nil {
		return err
	}

	// 已存在则不重复创建
	if count > 0 {
		log.Info("基础菜单已存在，跳过创建")
		return nil
	}

	// 创建系统管理菜单
	sysManage := model.Menu{
		Name:      "系统管理",
		ParentID:  0,
		Path:      "/system",
		Component: "Layout",
		Type:      0, // 目录
		Icon:      "system",
		OrderNum:  1,
		IsFrame:   false,
		IsHidden:  false,
		Enabled:   true,
		Perms:     "",
	}

	if err := c.db.Create(&sysManage).Error; err != nil {
		return err
	}

	// 创建用户管理菜单
	userManage := model.Menu{
		Name:      "用户管理",
		ParentID:  sysManage.ID,
		Path:      "user",
		Component: "system/user/index",
		Type:      1, // 菜单
		Icon:      "user",
		OrderNum:  1,
		IsFrame:   false,
		IsHidden:  false,
		Enabled:   true,
		Perms:     "system:user:list",
	}

	if err := c.db.Create(&userManage).Error; err != nil {
		return err
	}

	// 创建角色管理菜单
	roleManage := model.Menu{
		Name:      "角色管理",
		ParentID:  sysManage.ID,
		Path:      "role",
		Component: "system/role/index",
		Type:      1, // 菜单
		Icon:      "role",
		OrderNum:  2,
		IsFrame:   false,
		IsHidden:  false,
		Enabled:   true,
		Perms:     "system:role:list",
	}

	if err := c.db.Create(&roleManage).Error; err != nil {
		return err
	}

	// 创建菜单管理菜单
	menuManage := model.Menu{
		Name:      "菜单管理",
		ParentID:  sysManage.ID,
		Path:      "menu",
		Component: "system/menu/index",
		Type:      1, // 菜单
		Icon:      "menu",
		OrderNum:  3,
		IsFrame:   false,
		IsHidden:  false,
		Enabled:   true,
		Perms:     "system:menu:list",
	}

	if err := c.db.Create(&menuManage).Error; err != nil {
		return err
	}

	// 创建权限管理菜单
	permManage := model.Menu{
		Name:      "权限管理",
		ParentID:  sysManage.ID,
		Path:      "permission",
		Component: "system/permission/index",
		Type:      1, // 菜单
		Icon:      "permission",
		OrderNum:  4,
		IsFrame:   false,
		IsHidden:  false,
		Enabled:   true,
		Perms:     "system:permission:list",
	}

	if err := c.db.Create(&permManage).Error; err != nil {
		return err
	}

	// 为这些菜单创建相应的权限
	menus := []model.Menu{userManage, roleManage, menuManage, permManage}
	for _, menu := range menus {
		// 查看权限
		listPerm := model.Permission{
			Name:    menu.Name + "查看",
			Code:    menu.Perms,
			Remark:  "查看" + menu.Name + "列表",
			Type:    1, // 菜单权限
			Enabled: true,
			MenuID:  menu.ID,
		}

		if err := c.db.Create(&listPerm).Error; err != nil {
			return err
		}

		// 新增权限
		addPerm := model.Permission{
			Name:    menu.Name + "新增",
			Code:    menu.Perms[:len(menu.Perms)-4] + "add",
			Remark:  "新增" + menu.Name,
			Type:    2, // 操作权限
			Enabled: true,
			MenuID:  menu.ID,
		}

		if err := c.db.Create(&addPerm).Error; err != nil {
			return err
		}

		// 编辑权限
		editPerm := model.Permission{
			Name:    menu.Name + "编辑",
			Code:    menu.Perms[:len(menu.Perms)-4] + "edit",
			Remark:  "编辑" + menu.Name,
			Type:    2, // 操作权限
			Enabled: true,
			MenuID:  menu.ID,
		}

		if err := c.db.Create(&editPerm).Error; err != nil {
			return err
		}

		// 删除权限
		deletePerm := model.Permission{
			Name:    menu.Name + "删除",
			Code:    menu.Perms[:len(menu.Perms)-4] + "delete",
			Remark:  "删除" + menu.Name,
			Type:    2, // 操作权限
			Enabled: true,
			MenuID:  menu.ID,
		}

		if err := c.db.Create(&deletePerm).Error; err != nil {
			return err
		}
	}

	// 为管理员角色分配所有菜单
	var adminRole model.Role
	if err := c.db.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// 为超级管理员分配所有菜单
	allMenus := []model.Menu{sysManage, userManage, roleManage, menuManage, permManage}
	for _, menu := range allMenus {
		roleMenu := model.RoleMenu{
			RoleID: adminRole.ID,
			MenuID: menu.ID,
		}

		if err := c.db.Create(&roleMenu).Error; err != nil {
			return err
		}

		// 分配该菜单下的所有权限
		var permissions []model.Permission
		if err := c.db.Where("menu_id = ?", menu.ID).Find(&permissions).Error; err != nil {
			return err
		}

		for _, perm := range permissions {
			rolePerm := model.RolePermission{
				RoleID:       adminRole.ID,
				PermissionID: perm.ID,
			}

			if err := c.db.Create(&rolePerm).Error; err != nil {
				return err
			}

			// 添加到Casbin
			_, err := c.casbinComp.AddPermissionForRole("admin", perm.Code, "*")
			if err != nil {
				return err
			}
		}
	}

	log.Info("基础菜单创建成功")
	return nil
}
