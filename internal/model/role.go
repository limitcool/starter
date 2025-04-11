package model

// 角色实体
type Role struct {
	BaseModel

	Name        string `json:"name" gorm:"size:50;not null;unique;comment:角色名称"`
	Code        string `json:"code" gorm:"size:100;not null;unique;comment:角色编码"`
	Enabled     bool   `json:"enabled" gorm:"default:true;comment:是否启用"`
	Sort        int    `json:"sort" gorm:"default:0;comment:排序"`
	Description string `json:"description" gorm:"size:200;comment:角色描述"`
	Remark      string `json:"remark" gorm:"size:500;comment:备注"`

	// 关联
	Users         []*User       `json:"users" gorm:"many2many:sys_user_role;"`             // 关联的用户
	Permissions   []*Permission `json:"permissions" gorm:"many2many:sys_role_permission;"` // 关联的权限
	Menus         []*Menu       `json:"menus" gorm:"many2many:sys_role_menu;"`             // 关联的菜单
	UserIDs       []uint        `json:"user_ids" gorm:"-"`                                 // 用户ID列表，不映射到数据库
	PermissionIDs []uint        `json:"permission_ids" gorm:"-"`                           // 权限ID列表，不映射到数据库
	MenuIDs       []uint        `json:"menu_ids" gorm:"-"`                                 // 菜单ID列表，不映射到数据库
}

// 表名
func (Role) TableName() string {
	return "sys_role"
}

// 创建角色
func (r *Role) Create() error {
	return DB().Create(r).Error
}

// 更新角色
func (r *Role) Update() error {
	return DB().Model(&Role{}).Where("id = ?", r.ID).Updates(r).Error
}

// 删除角色
func (r *Role) Delete() error {
	return DB().Delete(&Role{}, r.ID).Error
}

// 根据ID获取角色
func (r *Role) GetByID(id uint) (*Role, error) {
	var role Role
	err := DB().Where("id = ?", id).First(&role).Error
	return &role, err
}

// 获取所有角色
func (r *Role) GetAll() ([]Role, error) {
	var roles []Role
	err := DB().Order("sort").Find(&roles).Error
	return roles, err
}

// 检查角色是否已分配给用户
func (r *Role) IsAssignedToUser(id uint) (bool, error) {
	var count int64
	err := DB().Model(&UserRole{}).Where("role_id = ?", id).Count(&count).Error
	return count > 0, err
}

// DeleteRoleMenus 删除角色的菜单关联
func (r *Role) DeleteRoleMenus(roleID uint) error {
	return DB().Where("role_id = ?", roleID).Delete(&RoleMenu{}).Error
}

// 角色菜单关联表
type RoleMenu struct {
	BaseModel

	RoleID uint `json:"role_id" gorm:"not null;comment:角色ID"`
	MenuID uint `json:"menu_id" gorm:"not null;comment:菜单ID"`
}

// 表名
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

// 批量创建角色菜单关联
func (rm *RoleMenu) BatchCreate(roleMenus []RoleMenu) error {
	return DB().Create(&roleMenus).Error
}

// 获取角色菜单ID列表
func (rm *RoleMenu) GetMenuIDsByRoleID(roleID uint) ([]uint, error) {
	var roleMenus []RoleMenu
	err := DB().Where("role_id = ?", roleID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	return menuIDs, nil
}

// 用户角色关联表
type UserRole struct {
	BaseModel

	UserID int64 `json:"user_id" gorm:"type:bigint;not null;comment:用户ID"`
	RoleID uint  `json:"role_id" gorm:"not null;comment:角色ID"`
}

// 表名
func (UserRole) TableName() string {
	return "sys_user_role"
}

// 批量创建用户角色关联
func (ur *UserRole) BatchCreate(userRoles []UserRole) error {
	return DB().Create(&userRoles).Error
}

// 删除用户的角色关联
func (ur *UserRole) DeleteByUserID(userID int64) error {
	return DB().Where("user_id = ?", userID).Delete(&UserRole{}).Error
}

// 获取用户的角色ID列表
func (ur *UserRole) GetRoleIDsByUserID(userID uint) ([]uint, error) {
	var userRoles []UserRole
	err := DB().Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	return roleIDs, nil
}
