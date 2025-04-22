package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/services"
	"github.com/spf13/cast"
)

func NewRoleController(roleService services.RoleServiceInterface, menuService services.MenuServiceInterface) *RoleController {
	return &RoleController{
		roleService: roleService,
		menuService: menuService,
	}
}

type RoleController struct {
	roleService services.RoleServiceInterface
	menuService services.MenuServiceInterface
}

// 创建角色
func (rc *RoleController) CreateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 使用控制器中的服务实例
	if err := rc.roleService.CreateRole(c.Request.Context(), &role); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 更新角色
func (rc *RoleController) UpdateRole(c *gin.Context) {
	id, err := cast.ToUint64E(c.Param("id"))
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	role.ID = uint(id)
	if err := rc.roleService.UpdateRole(c.Request.Context(), &role); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 删除角色
func (rc *RoleController) DeleteRole(c *gin.Context) {
	id, err := cast.ToUint64E(c.Param("id"))
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	if err := rc.roleService.DeleteRole(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 获取角色详情
func (rc *RoleController) GetRole(c *gin.Context) {
	id, err := cast.ToUint64E(c.Param("id"))
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 使用控制器中的服务实例
	role, err := rc.roleService.GetRoleByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	// 获取角色菜单ID
	roleMenus, err := rc.menuService.GetMenusByRoleID(c.Request.Context(), role.ID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 提取菜单ID
	menuIDs := make([]uint, 0)
	for _, menu := range roleMenus {
		menuIDs = append(menuIDs, menu.ID)
	}
	role.MenuIDs = menuIDs

	response.Success(c, role)
}

// 获取角色列表
func (rc *RoleController) GetRoles(c *gin.Context) {
	// 使用控制器中的服务实例
	roles, err := rc.roleService.GetRoles(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, roles)
}

// 为角色分配菜单
func (rc *RoleController) AssignMenuToRole(c *gin.Context) {
	var req struct {
		RoleID  uint   `json:"role_id" binding:"required"`
		MenuIDs []uint `json:"menu_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 使用控制器中的服务实例
	if err := rc.menuService.AssignMenuToRole(c.Request.Context(), req.RoleID, req.MenuIDs); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 为角色设置权限
func (rc *RoleController) SetRolePermission(c *gin.Context) {
	var req struct {
		RoleCode string `json:"role_code" binding:"required"`
		Object   string `json:"object" binding:"required"`
		Action   string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 使用控制器中的服务实例
	if err := rc.roleService.SetRolePermission(c.Request.Context(), req.RoleCode, req.Object, req.Action); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 删除角色权限
func (rc *RoleController) DeleteRolePermission(c *gin.Context) {
	var req struct {
		RoleCode string `json:"role_code" binding:"required"`
		Object   string `json:"object" binding:"required"`
		Action   string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 使用控制器中的服务实例
	if err := rc.roleService.DeleteRolePermission(c.Request.Context(), req.RoleCode, req.Object, req.Action); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}
