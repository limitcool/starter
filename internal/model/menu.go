package model

import (
	"errors"
	"sort"

	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/storage/sqldb"
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

// 创建菜单
func (m *Menu) Create() error {
	return sqldb.Instance().DB().Create(m).Error
}

// 更新菜单
func (m *Menu) Update() error {
	return sqldb.Instance().DB().Model(&Menu{}).Where("id = ?", m.ID).Updates(m).Error
}

// 删除菜单
func (m *Menu) Delete(id uint) error {
	// 检查是否有子菜单
	var count int64
	db := sqldb.Instance().DB()
	if err := db.Model(&Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该菜单下有子菜单，不能删除")
	}

	// 删除菜单
	return db.Delete(&Menu{}, id).Error
}

// 根据ID获取菜单
func (m *Menu) GetByID(id uint) (*Menu, error) {
	var menu Menu
	err := sqldb.Instance().DB().Where("id = ?", id).First(&menu).Error
	return &menu, err
}

// 获取所有菜单
func (m *Menu) GetAll() ([]*Menu, error) {
	var menus []*Menu
	err := sqldb.Instance().DB().Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return BuildMenuTreeOld(menus), nil
}

// 获取角色菜单
func (m *Menu) GetByRoleID(roleID uint) ([]*Menu, error) {
	// 查询角色关联的菜单ID
	var roleMenus []RoleMenu
	db := sqldb.Instance().DB()
	err := db.Where("role_id = ?", roleID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	// 如果没有菜单，返回空数组
	if len(menuIDs) == 0 {
		return []*Menu{}, nil
	}

	// 查询菜单
	var menus []*Menu
	err = db.Where("id IN ?", menuIDs).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return BuildMenuTreeOld(menus), nil
}

// 获取用户菜单
func (m *Menu) GetByUserID(userID uint) ([]*Menu, error) {
	// 1. 获取用户角色
	var userRoles []UserRole
	db := sqldb.Instance().DB()
	err := db.Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	// 如果用户没有角色，返回空数组
	if len(userRoles) == 0 {
		return []*Menu{}, nil
	}

	// 提取角色ID
	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	// 2. 获取角色关联的菜单ID
	var roleMenus []RoleMenu
	err = db.Where("role_id IN ?", roleIDs).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取菜单ID并去重
	menuIDMap := make(map[uint]bool)
	for _, rm := range roleMenus {
		menuIDMap[rm.MenuID] = true
	}

	if len(menuIDMap) == 0 {
		return []*Menu{}, nil
	}

	var menuIDs []uint
	for id := range menuIDMap {
		menuIDs = append(menuIDs, id)
	}

	// 3. 查询菜单信息
	var menus []*Menu
	err = db.Where("id IN ? AND status = ? AND type IN ?", menuIDs, 1, []int8{0, 1}).Order("order_num").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 4. 构建菜单树
	return BuildMenuTreeOld(menus), nil
}

// 获取用户菜单权限标识
func (m *Menu) GetPermsByUserRoles(roleIDs []uint) ([]string, error) {
	// 查询角色菜单
	var roleMenus []RoleMenu
	db := sqldb.Instance().DB()
	err := db.Where("role_id IN ?", roleIDs).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	// 提取菜单ID
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	if len(menuIDs) == 0 {
		return []string{}, nil
	}

	// 查询菜单权限标识
	var perms []string
	err = db.Model(&Menu{}).
		Where("id IN ? AND status = ? AND perms != ''", menuIDs, 1).
		Pluck("perms", &perms).Error
	if err != nil {
		return nil, err
	}

	return perms, nil
}

// BuildMenuTreeOld 构建菜单树(旧版本)
func BuildMenuTreeOld(menus []*Menu) []*Menu {
	// 创建一个map用于快速查找
	menuMap := make(map[uint]*Menu)
	for _, m := range menus {
		menuMap[m.ID] = m
	}

	var rootMenus []*Menu
	for _, m := range menus {
		if m.ParentID == 0 {
			// 顶级菜单
			rootMenus = append(rootMenus, m)
		} else {
			// 子菜单
			if parent, ok := menuMap[m.ParentID]; ok {
				if parent.Children == nil {
					parent.Children = []*Menu{}
				}
				parent.Children = append(parent.Children, m)
			}
		}
	}

	// 菜单排序
	for _, m := range menuMap {
		if len(m.Children) > 0 {
			sort.Slice(m.Children, func(i, j int) bool {
				return m.Children[i].OrderNum < m.Children[j].OrderNum
			})
		}
	}

	// 对根菜单排序
	sort.Slice(rootMenus, func(i, j int) bool {
		return rootMenus[i].OrderNum < rootMenus[j].OrderNum
	})

	return rootMenus
}
