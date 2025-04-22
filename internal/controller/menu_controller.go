package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/services"
	"github.com/spf13/cast"
)

// 创建菜单
func NewMenuController(menuService services.MenuServiceInterface) *MenuController {
	return &MenuController{
		menuService: menuService,
	}
}

type MenuController struct {
	menuService services.MenuServiceInterface
}

func (mc *MenuController) CreateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Error(c, err)
		return
	}

	// 使用控制器中的服务实例
	if err := mc.menuService.CreateMenu(c.Request.Context(), &menu); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 更新菜单
func (mc *MenuController) UpdateMenu(c *gin.Context) {
	id, err := cast.ToUint64E(c.Param("id"))
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Error(c, err)
		return
	}

	menu.ID = uint(id)
	if err := mc.menuService.UpdateMenu(c.Request.Context(), &menu); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 删除菜单
func (mc *MenuController) DeleteMenu(c *gin.Context) {
	id, err := cast.ToUint64E(c.Param("id"))
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	if err := mc.menuService.DeleteMenu(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 获取菜单详情
func (mc *MenuController) GetMenu(c *gin.Context) {
	id, err := cast.ToUint64E(c.Param("id"))
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 使用控制器中的服务实例
	menu, err := mc.menuService.GetMenuByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, menu)
}

// 获取菜单树
func (mc *MenuController) GetMenuTree(c *gin.Context) {
	menus, err := mc.menuService.GetMenuTree(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, menus)
}

// 获取用户菜单
func (mc *MenuController) GetUserMenus(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	// 使用控制器中的服务实例
	menus, err := mc.menuService.GetUserMenus(c.Request.Context(), cast.ToInt64(userID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, menus)
}

// 获取用户菜单权限标识
func (mc *MenuController) GetUserMenuPerms(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	// 使用控制器中的服务实例
	perms, err := mc.menuService.GetMenuPermsByUserID(c.Request.Context(), cast.ToUint(userID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, perms)
}
