package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/permission"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	permissionService *permission.Service
	roleRepo          *model.RoleRepo
	permRepo          *model.PermissionRepo
	menuRepo          *model.MenuRepo
	userRepo          *model.UserRepo
}

// NewPermissionHandler 创建权限处理器
func NewPermissionHandler(
	permissionService *permission.Service,
	roleRepo *model.RoleRepo,
	permRepo *model.PermissionRepo,
	menuRepo *model.MenuRepo,
	userRepo *model.UserRepo,
) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
		roleRepo:          roleRepo,
		permRepo:          permRepo,
		menuRepo:          menuRepo,
		userRepo:          userRepo,
	}
}

// GetUserMenus 获取用户菜单
func (h *PermissionHandler) GetUserMenus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errorx.ErrUnauthorized)
		return
	}

	uid, ok := userID.(int64)
	if !ok {
		response.Error(c, errorx.ErrUnauthorized)
		return
	}

	menus, err := h.permissionService.GetUserMenus(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, menus)
}

// GetUserPermissions 获取用户权限
func (h *PermissionHandler) GetUserPermissions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errorx.ErrUnauthorized)
		return
	}

	uid, ok := userID.(int64)
	if !ok {
		response.Error(c, errorx.ErrUnauthorized)
		return
	}

	permissions, err := h.permissionService.GetUserPermissions(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, permissions)
}

// CheckPermission 检查权限
func (h *PermissionHandler) CheckPermission(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errorx.ErrUnauthorized)
		return
	}

	uid, ok := userID.(int64)
	if !ok {
		response.Error(c, errorx.ErrUnauthorized)
		return
	}

	var req dto.CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	hasPermission, err := h.permissionService.CheckPermission(c.Request.Context(), uid, req.Resource, req.Action)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.CheckPermissionResponse{
		HasPermission: hasPermission,
	})
}

// ListRoles 获取角色列表
func (h *PermissionHandler) ListRoles(c *gin.Context) {
	var req dto.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	// 构建查询选项
	opts := &model.QueryOptions{}
	if req.Keyword != "" {
		opts.Condition = "name LIKE ? OR display_name LIKE ?"
		opts.Args = []any{"%" + req.Keyword + "%", "%" + req.Keyword + "%"}
	}
	if req.Status != nil {
		if opts.Condition != "" {
			opts.Condition += " AND status = ?"
		} else {
			opts.Condition = "status = ?"
		}
		opts.Args = append(opts.Args, *req.Status)
	}

	roles, err := h.roleRepo.List(c.Request.Context(), req.Page, req.PageSize, opts)
	if err != nil {
		response.Error(c, err)
		return
	}

	total, err := h.roleRepo.Count(c.Request.Context(), opts)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPagination(c, roles, total, req.Page, req.PageSize)
}

// CreateRole 创建角色
func (h *PermissionHandler) CreateRole(c *gin.Context) {
	var req dto.RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	role := &model.Role{
		Name:        req.Name,
		Key:         req.Key,
		Description: req.Description,
		Status:      uint8(req.Status),
	}

	if err := h.roleRepo.Create(c.Request.Context(), role); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
func (h *PermissionHandler) UpdateRole(c *gin.Context) {
	var req dto.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	role, err := h.roleRepo.Get(c.Request.Context(), req.ID, nil)
	if err != nil {
		response.Error(c, err)
		return
	}

	role.Name = req.Name
	role.Key = req.Key
	role.Description = req.Description
	role.Status = uint8(req.Status)

	if err := h.roleRepo.Update(c.Request.Context(), role); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
func (h *PermissionHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("无效的角色ID"))
		return
	}

	if err := h.roleRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, struct{}{})
}

// AssignUserRoles 分配用户角色
func (h *PermissionHandler) AssignUserRoles(c *gin.Context) {
	var req dto.AssignUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	if err := h.permissionService.AssignRolesToUser(c.Request.Context(), req.UserID, req.RoleKeys); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, struct{}{})
}

// AssignRolePermissions 分配角色权限
func (h *PermissionHandler) AssignRolePermissions(c *gin.Context) {
	var req dto.AssignRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	if err := h.permissionService.AssignPermissionsToRole(c.Request.Context(), req.RoleID, req.PermissionKeys); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, struct{}{})
}

// ListPermissions 获取权限列表
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	var req dto.PermissionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	// 构建查询选项
	opts := &model.QueryOptions{}
	conditions := []string{}
	args := []any{}

	if req.Keyword != "" {
		conditions = append(conditions, "(name LIKE ? OR display_name LIKE ?)")
		args = append(args, "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *req.Status)
	}
	if req.Resource != "" {
		conditions = append(conditions, "resource = ?")
		args = append(args, req.Resource)
	}
	if req.Action != "" {
		conditions = append(conditions, "action = ?")
		args = append(args, req.Action)
	}

	if len(conditions) > 0 {
		opts.Condition = conditions[0]
		for i := 1; i < len(conditions); i++ {
			opts.Condition += " AND " + conditions[i]
		}
		opts.Args = args
	}

	permissions, err := h.permRepo.List(c.Request.Context(), req.Page, req.PageSize, opts)
	if err != nil {
		response.Error(c, err)
		return
	}

	total, err := h.permRepo.Count(c.Request.Context(), opts)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPagination(c, permissions, total, req.Page, req.PageSize)
}

// ListMenus 获取菜单列表
func (h *PermissionHandler) ListMenus(c *gin.Context) {
	var req dto.MenuListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	// 构建查询选项
	opts := &model.QueryOptions{}
	conditions := []string{}
	args := []any{}

	if req.Keyword != "" {
		conditions = append(conditions, "(name LIKE ? OR title LIKE ?)")
		args = append(args, "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *req.Status)
	}
	if req.Type != nil {
		conditions = append(conditions, "type = ?")
		args = append(args, *req.Type)
	}
	if req.ParentID != nil {
		conditions = append(conditions, "parent_id = ?")
		args = append(args, *req.ParentID)
	}

	if len(conditions) > 0 {
		opts.Condition = conditions[0]
		for i := 1; i < len(conditions); i++ {
			opts.Condition += " AND " + conditions[i]
		}
		opts.Args = args
	}

	menus, err := h.menuRepo.List(c.Request.Context(), 1, 1000, opts)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 构建菜单树
	menuTree := h.menuRepo.BuildMenuTree(menus, 0)

	response.Success(c, menuTree)
}
