package model

import (
	"github.com/limitcool/starter/internal/pkg/enum"
)

// Menu 菜单实体
type Menu struct {
	BaseModel

	Name      string        `json:"name" gorm:"size:50;not null;comment:菜单名称"`
	ParentID  uint          `json:"parent_id" gorm:"default:0;index;comment:父菜单ID"`
	Path      string        `json:"path" gorm:"size:100;comment:前端路由路径"`
	Component string        `json:"component" gorm:"size:100;comment:前端组件路径"`
	Perms     string        `json:"perms" gorm:"size:100;index;comment:权限标识"`
	Type      enum.MenuType `json:"type" gorm:"default:0;index;comment:菜单类型(0:目录,1:菜单,2:按钮)"`
	Icon      string        `json:"icon" gorm:"size:50;comment:图标"`
	OrderNum  int           `json:"order_num" gorm:"default:0;comment:排序号"`
	IsFrame   bool          `json:"is_frame" gorm:"default:false;comment:是否为外链"`
	IsHidden  bool          `json:"is_hidden" gorm:"default:false;comment:是否隐藏"`
	Enabled   bool          `json:"enabled" gorm:"default:true;comment:是否启用"`
	Remark    string        `json:"remark" gorm:"size:500;comment:备注"`
	Redirect  string        `json:"redirect" gorm:"size:100;comment:重定向路径"`
	Title     string        `json:"title" gorm:"size:50;comment:标题"`
	KeepAlive bool          `json:"keep_alive" gorm:"default:false;comment:是否缓存"`

	Children []*Menu `json:"children" gorm:"-"` // 子菜单，不映射到数据库

	// 关联
	Buttons []*MenuButton `json:"buttons" gorm:"foreignKey:MenuID"`    // 按钮
	APIs    []*API        `json:"apis" gorm:"many2many:sys_menu_api;"` // 关联的API

	// 权限关联
	Permissions   []Permission `json:"permissions" gorm:"foreignKey:MenuID"` // 关联的权限列表
	PermissionIDs []uint       `json:"permission_ids" gorm:"-"`              // 权限ID列表，不映射到数据库
}

// MenuButton 菜单按钮
type MenuButton struct {
	BaseModel

	MenuID     uint   `json:"menu_id" gorm:"not null;index;comment:菜单ID"`
	Name       string `json:"name" gorm:"size:50;not null;comment:按钮名称"`
	Permission string `json:"permission" gorm:"size:100;not null;index;comment:权限标识"`
	Icon       string `json:"icon" gorm:"size:50;comment:图标"`
	OrderNum   int    `json:"order_num" gorm:"default:0;comment:排序号"`
	Enabled    bool   `json:"enabled" gorm:"default:true;comment:是否启用"`
	Remark     string `json:"remark" gorm:"size:500;comment:备注"`
}

// API 接口实体
type API struct {
	BaseModel

	Path        string `json:"path" gorm:"size:100;not null;comment:接口路径"`
	Method      string `json:"method" gorm:"size:10;not null;comment:请求方法"`
	Name        string `json:"name" gorm:"size:100;comment:接口名称"`
	Description string `json:"description" gorm:"size:200;comment:接口描述"`
	Group       string `json:"group" gorm:"size:50;comment:接口分组"`
	Enabled     bool   `json:"enabled" gorm:"default:true;comment:是否启用"`
	Remark      string `json:"remark" gorm:"size:500;comment:备注"`
}

// MenuAPI 菜单与API关联
type MenuAPI struct {
	BaseModel

	MenuID uint `json:"menu_id" gorm:"not null;index;comment:菜单ID"`
	APIID  uint `json:"api_id" gorm:"not null;index;comment:接口ID"`
}

// 表名
func (Menu) TableName() string {
	return "sys_menu"
}

// 表名
func (MenuButton) TableName() string {
	return "sys_menu_button"
}

// 表名
func (API) TableName() string {
	return "sys_api"
}

// 表名
func (MenuAPI) TableName() string {
	return "sys_menu_api"
}

// 以下方法已移动到 repository/menu_repo.go
// Create
// Update
// Delete
// GetByID
// GetAll
// GetByRoleID
// GetByUserID
// GetPermsByUserRoles

// 已移除 BuildMenuTreeOld 方法到 repository/menu_repo.go
