package model

import (
	"github.com/limitcool/starter/internal/pkg/enum"
)

// 菜单实体
type Menu struct {
	BaseModel

	Name      string        `json:"name" gorm:"size:50;not null;comment:菜单名称"`
	ParentID  uint          `json:"parent_id" gorm:"default:0;comment:父菜单ID"`
	Path      string        `json:"path" gorm:"size:100;comment:前端路由路径"`
	Component string        `json:"component" gorm:"size:100;comment:前端组件路径"`
	Perms     string        `json:"perms" gorm:"size:100;comment:权限标识"`
	Type      enum.MenuType `json:"type" gorm:"default:0;comment:菜单类型(0:目录,1:菜单,2:按钮)"`
	Icon      string        `json:"icon" gorm:"size:50;comment:图标"`
	OrderNum  int           `json:"order_num" gorm:"default:0;comment:排序号"`
	IsFrame   bool          `json:"is_frame" gorm:"default:false;comment:是否为外链"`
	IsHidden  bool          `json:"is_hidden" gorm:"default:false;comment:是否隐藏"`
	Enabled   bool          `json:"enabled" gorm:"default:true;comment:是否启用"`
	Remark    string        `json:"remark" gorm:"size:500;comment:备注"`

	Children []*Menu `json:"children" gorm:"-"` // 子菜单，不映射到数据库

	// 权限关联
	Permissions   []Permission `json:"permissions" gorm:"foreignKey:MenuID"` // 关联的权限列表
	PermissionIDs []uint       `json:"permission_ids" gorm:"-"`              // 权限ID列表，不映射到数据库
}

// 表名
func (Menu) TableName() string {
	return "sys_menu"
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
