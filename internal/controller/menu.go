package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/pkg/response"
)

// 创建菜单
func CreateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	menuService := services.NewMenuService(global.DB)
	if err := menuService.CreateMenu(&menu); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// 更新菜单
func UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	menu.ID = uint(id)
	menuService := services.NewMenuService(global.DB)
	if err := menuService.UpdateMenu(&menu); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// 删除菜单
func DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	menuService := services.NewMenuService(global.DB)
	if err := menuService.DeleteMenu(uint(id)); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// 获取菜单详情
func GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	menuService := services.NewMenuService(global.DB)
	menu, err := menuService.GetMenuByID(uint(id))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, menu)
}

// 获取菜单树
func GetMenuTree(c *gin.Context) {
	menuService := services.NewMenuService(global.DB)
	menus, err := menuService.GetMenuTree()
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, menus)
}

// 获取用户菜单
func GetUserMenus(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	menuService := services.NewMenuService(global.DB)
	menus, err := menuService.GetUserMenus(uint(userID.(float64)))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, menus)
}

// 获取用户菜单权限标识
func GetUserMenuPerms(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	menuService := services.NewMenuService(global.DB)
	perms, err := menuService.GetMenuPermsByUserID(uint(userID.(float64)))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, perms)
}
