package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/services"
)

// 创建菜单
func NewMenuController() *MenuController {
	return &MenuController{}
}

type MenuController struct {
}

func (mc *MenuController) CreateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Error(c, err)
		return
	}

	menuService := services.NewMenuService()
	if err := menuService.CreateMenu(&menu); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 更新菜单
func (mc *MenuController) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
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
	menuService := services.NewMenuService()
	if err := menuService.UpdateMenu(&menu); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 删除菜单
func (mc *MenuController) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	menuService := services.NewMenuService()
	if err := menuService.DeleteMenu(uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// 获取菜单详情
func (mc *MenuController) GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	menuService := services.NewMenuService()
	menu, err := menuService.GetMenuByID(uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, menu)
}

// 获取菜单树
func (mc *MenuController) GetMenuTree(c *gin.Context) {
	menuService := services.NewMenuService()
	menus, err := menuService.GetMenuTree()
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

	menuService := services.NewMenuService()
	menus, err := menuService.GetUserMenus(uint(userID.(float64)))
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

	menuService := services.NewMenuService()
	perms, err := menuService.GetMenuPermsByUserID(uint(userID.(float64)))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, perms)
}
