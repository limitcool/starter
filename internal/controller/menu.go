package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/internal/pkg/apiresponse"
)

// 创建菜单
func CreateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		apiresponse.ParamError(c, err.Error())
		return
	}

	db := services.Instance().GetDB()
	menuService := services.NewMenuService(db)
	if err := menuService.CreateMenu(&menu); err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success[any](c, nil)
}

// 更新菜单
func UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiresponse.ParamError(c, "无效的菜单ID")
		return
	}

	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		apiresponse.ParamError(c, err.Error())
		return
	}

	menu.ID = uint(id)
	db := services.Instance().GetDB()
	menuService := services.NewMenuService(db)
	if err := menuService.UpdateMenu(&menu); err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success[any](c, nil)
}

// 删除菜单
func DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiresponse.ParamError(c, "无效的菜单ID")
		return
	}

	db := services.Instance().GetDB()
	menuService := services.NewMenuService(db)
	if err := menuService.DeleteMenu(uint(id)); err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success[any](c, nil)
}

// 获取菜单详情
func GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiresponse.ParamError(c, "无效的菜单ID")
		return
	}

	db := services.Instance().GetDB()
	menuService := services.NewMenuService(db)
	menu, err := menuService.GetMenuByID(uint(id))
	if err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success(c, menu)
}

// 获取菜单树
func GetMenuTree(c *gin.Context) {
	db := services.Instance().GetDB()
	menuService := services.NewMenuService(db)
	menus, err := menuService.GetMenuTree()
	if err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success(c, menus)
}

// 获取用户菜单
func GetUserMenus(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		apiresponse.Unauthorized(c, "未登录")
		return
	}

	db := services.Instance().GetDB()
	menuService := services.NewMenuService(db)
	menus, err := menuService.GetUserMenus(uint(userID.(float64)))
	if err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success(c, menus)
}

// 获取用户菜单权限标识
func GetUserMenuPerms(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		apiresponse.Unauthorized(c, "未登录")
		return
	}

	db := services.Instance().GetDB()
	menuService := services.NewMenuService(db)
	perms, err := menuService.GetMenuPermsByUserID(uint(userID.(float64)))
	if err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success(c, perms)
}
