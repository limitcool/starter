package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// MenuCRUDHandler 菜单CRUD处理器
type MenuCRUDHandler struct {
	menuRepo *model.MenuRepo
}

// NewMenuCRUDHandler 创建菜单CRUD处理器
func NewMenuCRUDHandler(menuRepo *model.MenuRepo) *MenuCRUDHandler {
	return &MenuCRUDHandler{
		menuRepo: menuRepo,
	}
}

// GetMenus 获取菜单列表
func (h *MenuCRUDHandler) GetMenus(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	parentID, _ := strconv.Atoi(c.DefaultQuery("parent_id", "0"))
	platform := c.DefaultQuery("platform", "admin")

	// 构建查询条件
	conditions := []string{"platform = ?"}
	args := []any{platform}

	if parentID > 0 {
		conditions = append(conditions, "parent_id = ?")
		args = append(args, parentID)
	}

	var opts *model.QueryOptions
	if len(conditions) > 0 {
		condition := conditions[0]
		for i := 1; i < len(conditions); i++ {
			condition += " AND " + conditions[i]
		}
		opts = &model.QueryOptions{
			Condition: condition,
			Args:      args,
		}
	}

	menus, total, err := h.menuRepo.ListWithPagination(c.Request.Context(), page, pageSize, opts)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取菜单列表失败", "error", err)
		response.Error(c, err)
		return
	}

	response.SuccessWithPagination(c, menus, total, page, pageSize)
}

// GetMenuTree 获取菜单树
func (h *MenuCRUDHandler) GetMenuTree(c *gin.Context) {
	platform := c.DefaultQuery("platform", "admin")

	opts := &model.QueryOptions{
		Condition: "platform = ?",
		Args:      []any{platform},
	}

	menus, err := h.menuRepo.GetAll(c.Request.Context(), opts)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取菜单树失败", "error", err)
		response.Error(c, err)
		return
	}

	// 构建菜单树
	tree := h.buildMenuTree(menus, 0)
	response.Success(c, tree)
}

// GetMenu 获取菜单详情
func (h *MenuCRUDHandler) GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("菜单ID格式错误"))
		return
	}

	menu, err := h.menuRepo.Get(c.Request.Context(), uint(id), nil)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取菜单详情失败", "id", id, "error", err)
		response.Error(c, err)
		return
	}

	response.Success(c, menu)
}

// CreateMenu 创建菜单
func (h *MenuCRUDHandler) CreateMenu(c *gin.Context) {
	var req dto.MenuCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	menu := &model.Menu{
		ParentID:      req.ParentID,
		Name:          req.Name,
		Path:          req.Path,
		Component:     req.Component,
		Icon:          req.Icon,
		SortOrder:     req.SortOrder,
		IsVisible:     req.IsVisible,
		PermissionKey: req.PermissionKey,
		Platform:      req.Platform,
	}

	if err := h.menuRepo.Create(c.Request.Context(), menu); err != nil {
		logger.ErrorContext(c.Request.Context(), "创建菜单失败", "menu", menu, "error", err)
		response.Error(c, err)
		return
	}

	logger.InfoContext(c.Request.Context(), "创建菜单成功", "menu_id", menu.ID, "name", menu.Name)
	response.Success(c, menu)
}

// UpdateMenu 更新菜单
func (h *MenuCRUDHandler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("菜单ID格式错误"))
		return
	}

	var req dto.MenuUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	menu, err := h.menuRepo.Get(c.Request.Context(), uint(id), nil)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取菜单失败", "id", id, "error", err)
		response.Error(c, err)
		return
	}

	menu.ParentID = req.ParentID
	menu.Name = req.Name
	menu.Path = req.Path
	menu.Component = req.Component
	menu.Icon = req.Icon
	menu.SortOrder = req.SortOrder
	menu.IsVisible = req.IsVisible
	menu.PermissionKey = req.PermissionKey
	menu.Platform = req.Platform

	if err := h.menuRepo.Update(c.Request.Context(), menu); err != nil {
		logger.ErrorContext(c.Request.Context(), "更新菜单失败", "id", id, "error", err)
		response.Error(c, err)
		return
	}

	logger.InfoContext(c.Request.Context(), "更新菜单成功", "menu_id", menu.ID, "name", menu.Name)
	response.Success(c, menu)
}

// DeleteMenu 删除菜单
func (h *MenuCRUDHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("菜单ID格式错误"))
		return
	}

	// 检查是否有子菜单
	hasChildren, err := h.menuRepo.HasChildren(c.Request.Context(), uint(id))
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "检查子菜单失败", "id", id, "error", err)
		response.Error(c, err)
		return
	}

	if hasChildren {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("该菜单下还有子菜单，无法删除"))
		return
	}

	if err := h.menuRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		logger.ErrorContext(c.Request.Context(), "删除菜单失败", "id", id, "error", err)
		response.Error(c, err)
		return
	}

	logger.InfoContext(c.Request.Context(), "删除菜单成功", "menu_id", id)
	response.Success(c, struct{}{})
}

// UpdateMenuSort 更新菜单排序
func (h *MenuCRUDHandler) UpdateMenuSort(c *gin.Context) {
	var req dto.MenuSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	if err := h.menuRepo.UpdateSort(c.Request.Context(), req.MenuSorts); err != nil {
		logger.ErrorContext(c.Request.Context(), "更新菜单排序失败", "error", err)
		response.Error(c, err)
		return
	}

	logger.InfoContext(c.Request.Context(), "更新菜单排序成功")
	response.Success(c, struct{}{})
}

// buildMenuTree 构建菜单树
func (h *MenuCRUDHandler) buildMenuTree(menus []model.Menu, parentID uint) []dto.MenuTreeNode {
	var tree []dto.MenuTreeNode

	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := dto.MenuTreeNode{
				ID:            menu.ID,
				ParentID:      menu.ParentID,
				Name:          menu.Name,
				Path:          menu.Path,
				Component:     menu.Component,
				Icon:          menu.Icon,
				SortOrder:     menu.SortOrder,
				IsVisible:     menu.IsVisible,
				PermissionKey: menu.PermissionKey,
				Platform:      menu.Platform,
				Children:      h.buildMenuTree(menus, menu.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}
