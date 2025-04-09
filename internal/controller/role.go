package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/pkg/apiresponse"
	"github.com/limitcool/starter/pkg/code"
)

// 创建角色
func CreateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		apiresponse.HandleError(c, code.NewErrCodeMsg(code.InvalidParams, err.Error()))
		return
	}

	roleService := services.NewRoleService(global.DB)
	if err := roleService.CreateRole(&role); err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[any](c, nil)
}

// 更新角色
func UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiresponse.HandleError(c, code.NewErrCodeMsg(code.InvalidParams, "无效的角色ID"))
		return
	}

	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		apiresponse.HandleError(c, code.NewErrCodeMsg(code.InvalidParams, err.Error()))
		return
	}

	role.ID = uint(id)
	roleService := services.NewRoleService(global.DB)
	if err := roleService.UpdateRole(&role); err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[any](c, nil)
}

// 删除角色
func DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiresponse.HandleError(c, code.NewErrCodeMsg(code.InvalidParams, "无效的角色ID"))
		return
	}

	roleService := services.NewRoleService(global.DB)
	if err := roleService.DeleteRole(uint(id)); err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[any](c, nil)
}

// 获取角色详情
func GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiresponse.HandleError(c, code.NewErrCodeMsg(code.InvalidParams, "无效的角色ID"))
		return
	}

	roleService := services.NewRoleService(global.DB)
	role, err := roleService.GetRoleByID(uint(id))
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	// 获取角色菜单ID
	menuService := services.NewMenuService(global.DB)
	roleMenus, err := menuService.GetMenusByRoleID(role.ID)
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	// 提取菜单ID
	menuIDs := make([]uint, 0)
	for _, menu := range roleMenus {
		menuIDs = append(menuIDs, menu.ID)
	}
	role.MenuIDs = menuIDs

	apiresponse.Success(c, role)
}

// 获取角色列表
func GetRoles(c *gin.Context) {
	roleService := services.NewRoleService(global.DB)
	roles, err := roleService.GetRoles()
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success(c, roles)
}

// 为角色分配菜单
func AssignMenuToRole(c *gin.Context) {
	var req struct {
		RoleID  uint   `json:"role_id" binding:"required"`
		MenuIDs []uint `json:"menu_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.HandleError(c, code.NewErrCodeMsg(code.InvalidParams, err.Error()))
		return
	}

	menuService := services.NewMenuService(global.DB)
	if err := menuService.AssignMenuToRole(req.RoleID, req.MenuIDs); err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[any](c, nil)
}

// 为角色设置权限
func SetRolePermission(c *gin.Context) {
	var req struct {
		RoleCode string `json:"role_code" binding:"required"`
		Object   string `json:"object" binding:"required"`
		Action   string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.HandleError(c, code.NewErrCodeMsg(code.InvalidParams, err.Error()))
		return
	}

	roleService := services.NewRoleService(global.DB)
	if err := roleService.SetRolePermission(req.RoleCode, req.Object, req.Action); err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[any](c, nil)
}

// 删除角色权限
func DeleteRolePermission(c *gin.Context) {
	var req struct {
		RoleCode string `json:"role_code" binding:"required"`
		Object   string `json:"object" binding:"required"`
		Action   string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.HandleError(c, code.NewErrCodeMsg(code.InvalidParams, err.Error()))
		return
	}

	roleService := services.NewRoleService(global.DB)
	if err := roleService.DeleteRolePermission(req.RoleCode, req.Object, req.Action); err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[any](c, nil)
}
