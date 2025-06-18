package dto

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required,max=50" label:"角色名称"`
	Key         string `json:"key" binding:"required,max=50" label:"角色标识"`
	Description string `json:"description" binding:"max=255" label:"角色描述"`
	Status      int    `json:"status" binding:"oneof=0 1" label:"状态"`
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	ID          uint   `json:"id" binding:"required" label:"角色ID"`
	Name        string `json:"name" binding:"required,max=50" label:"角色名称"`
	Key         string `json:"key" binding:"required,max=50" label:"角色标识"`
	Description string `json:"description" binding:"max=255" label:"角色描述"`
	Status      int    `json:"status" binding:"oneof=0 1" label:"状态"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// RoleListRequest 角色列表请求
type RoleListRequest struct {
	PageRequest
	Keyword string `form:"keyword" label:"关键字"`
	Status  *int   `form:"status" binding:"omitempty,oneof=0 1" label:"状态"`
}

// AssignRolePermissionsRequest 分配角色权限请求
type AssignRolePermissionsRequest struct {
	RoleID         uint     `json:"role_id" binding:"required" label:"角色ID"`
	PermissionKeys []string `json:"permission_keys" binding:"required" label:"权限标识列表"`
}

// PermissionCreateRequest 创建权限请求
type PermissionCreateRequest struct {
	ParentID uint   `json:"parent_id" label:"父权限ID"`
	Name     string `json:"name" binding:"required,max=50" label:"权限名称"`
	Key      string `json:"key" binding:"required,max=100" label:"权限标识"`
	Type     string `json:"type" binding:"required,oneof=MENU BUTTON API" label:"权限类型"`
}

// PermissionUpdateRequest 更新权限请求
type PermissionUpdateRequest struct {
	ID       uint   `json:"id" binding:"required" label:"权限ID"`
	ParentID uint   `json:"parent_id" label:"父权限ID"`
	Name     string `json:"name" binding:"required,max=50" label:"权限名称"`
	Key      string `json:"key" binding:"required,max=100" label:"权限标识"`
	Type     string `json:"type" binding:"required,oneof=MENU BUTTON API" label:"权限类型"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID        uint   `json:"id"`
	ParentID  uint   `json:"parent_id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// PermissionTreeNode 权限树节点
type PermissionTreeNode struct {
	ID       uint                 `json:"id"`
	ParentID uint                 `json:"parent_id"`
	Name     string               `json:"name"`
	Key      string               `json:"key"`
	Type     string               `json:"type"`
	Children []PermissionTreeNode `json:"children"`
}

// AssignRoleMenusRequest 分配角色菜单请求
type AssignRoleMenusRequest struct {
	RoleID  uint   `json:"role_id" binding:"required" label:"角色ID"`
	MenuIDs []uint `json:"menu_ids" binding:"required" label:"菜单ID列表"`
}

// AssignUserRolesRequest 分配用户角色请求
type AssignUserRolesRequest struct {
	UserID   int64    `json:"user_id" binding:"required" label:"用户ID"`
	RoleKeys []string `json:"role_keys" binding:"required" label:"角色标识列表"`
}

// PermissionListRequest 权限列表请求
type PermissionListRequest struct {
	PageRequest
	Keyword  string `form:"keyword" label:"关键字"`
	Status   *int   `form:"status" binding:"omitempty,oneof=0 1" label:"状态"`
	Resource string `form:"resource" label:"资源标识"`
	Action   string `form:"action" label:"操作标识"`
}

// MenuCreateRequest 创建菜单请求
type MenuCreateRequest struct {
	ParentID      uint   `json:"parent_id" label:"父级ID"`
	Name          string `json:"name" binding:"required,max=50" label:"菜单名称"`
	Path          string `json:"path" binding:"max=255" label:"路由路径"`
	Component     string `json:"component" binding:"max=255" label:"组件路径"`
	Icon          string `json:"icon" binding:"max=50" label:"菜单图标"`
	SortOrder     int    `json:"sort_order" label:"排序"`
	IsVisible     bool   `json:"is_visible" label:"是否可见"`
	PermissionKey string `json:"permission_key" binding:"max=100" label:"权限标识"`
	Platform      string `json:"platform" binding:"required,oneof=admin coach_mp" label:"所属平台"`
}

// MenuUpdateRequest 更新菜单请求
type MenuUpdateRequest struct {
	ID            uint   `json:"id" binding:"required" label:"菜单ID"`
	ParentID      uint   `json:"parent_id" label:"父级ID"`
	Name          string `json:"name" binding:"required,max=50" label:"菜单名称"`
	Path          string `json:"path" binding:"max=255" label:"路由路径"`
	Component     string `json:"component" binding:"max=255" label:"组件路径"`
	Icon          string `json:"icon" binding:"max=50" label:"菜单图标"`
	SortOrder     int    `json:"sort_order" label:"排序"`
	IsVisible     bool   `json:"is_visible" label:"是否可见"`
	PermissionKey string `json:"permission_key" binding:"max=100" label:"权限标识"`
	Platform      string `json:"platform" binding:"required,oneof=admin coach_mp" label:"所属平台"`
}

// MenuResponse 菜单响应
type MenuResponse struct {
	ID            uint           `json:"id"`
	ParentID      uint           `json:"parent_id"`
	Name          string         `json:"name"`
	Path          string         `json:"path"`
	Component     string         `json:"component"`
	Icon          string         `json:"icon"`
	SortOrder     int            `json:"sort_order"`
	IsVisible     bool           `json:"is_visible"`
	PermissionKey string         `json:"permission_key"`
	Platform      string         `json:"platform"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
	Children      []MenuResponse `json:"children,omitempty"`
}

// MenuTreeNode 菜单树节点
type MenuTreeNode struct {
	ID            uint           `json:"id"`
	ParentID      uint           `json:"parent_id"`
	Name          string         `json:"name"`
	Path          string         `json:"path"`
	Component     string         `json:"component"`
	Icon          string         `json:"icon"`
	SortOrder     int            `json:"sort_order"`
	IsVisible     bool           `json:"is_visible"`
	PermissionKey string         `json:"permission_key"`
	Platform      string         `json:"platform"`
	Children      []MenuTreeNode `json:"children"`
}

// MenuSortRequest 菜单排序请求
type MenuSortRequest struct {
	MenuSorts []MenuSort `json:"menu_sorts" binding:"required" label:"菜单排序列表"`
}

// MenuSort 菜单排序
type MenuSort struct {
	ID        uint `json:"id" binding:"required" label:"菜单ID"`
	SortOrder int  `json:"sort_order" label:"排序"`
}

// MenuListRequest 菜单列表请求
type MenuListRequest struct {
	Keyword  string `form:"keyword" label:"关键字"`
	Status   *int   `form:"status" binding:"omitempty,oneof=0 1" label:"状态"`
	Type     *int   `form:"type" binding:"omitempty,oneof=1 2" label:"类型"`
	ParentID *uint  `form:"parent_id" label:"父级ID"`
}

// AssignMenuPermissionsRequest 分配菜单权限请求
type AssignMenuPermissionsRequest struct {
	MenuID        uint   `json:"menu_id" binding:"required" label:"菜单ID"`
	PermissionIDs []uint `json:"permission_ids" binding:"required" label:"权限ID列表"`
}

// UserPermissionInfoResponse 用户权限信息响应
type UserPermissionInfoResponse struct {
	UserID      int64                `json:"user_id"`
	Username    string               `json:"username"`
	IsAdmin     bool                 `json:"is_admin"`
	Roles       []RoleResponse       `json:"roles"`
	Permissions []PermissionResponse `json:"permissions"`
	Menus       []MenuResponse       `json:"menus"`
}

// CheckPermissionRequest 检查权限请求
type CheckPermissionRequest struct {
	Resource string `json:"resource" binding:"required" label:"资源标识"`
	Action   string `json:"action" binding:"required" label:"操作标识"`
}

// CheckPermissionResponse 检查权限响应
type CheckPermissionResponse struct {
	HasPermission bool `json:"has_permission"`
}
