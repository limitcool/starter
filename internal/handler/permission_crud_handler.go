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

// PermissionCRUDHandler 权限CRUD处理器
type PermissionCRUDHandler struct {
	permissionRepo *model.PermissionRepo
}

// NewPermissionCRUDHandler 创建权限CRUD处理器
func NewPermissionCRUDHandler(permissionRepo *model.PermissionRepo) *PermissionCRUDHandler {
	return &PermissionCRUDHandler{
		permissionRepo: permissionRepo,
	}
}

// GetPermissions 获取权限列表
func (h *PermissionCRUDHandler) GetPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	parentID, _ := strconv.Atoi(c.DefaultQuery("parent_id", "0"))

	var opts *model.QueryOptions
	if parentID > 0 {
		opts = &model.QueryOptions{
			Condition: "parent_id = ?",
			Args:      []any{parentID},
		}
	}

	permissions, total, err := h.permissionRepo.ListWithPagination(c.Request.Context(), page, pageSize, opts)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取权限列表失败", "error", err)
		response.Error(c, err)
		return
	}

	response.SuccessWithPagination(c, permissions, total, page, pageSize)
}

// GetPermissionTree 获取权限树
func (h *PermissionCRUDHandler) GetPermissionTree(c *gin.Context) {
	permissions, err := h.permissionRepo.GetAll(c.Request.Context())
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取权限树失败", "error", err)
		response.Error(c, err)
		return
	}

	// 构建权限树
	tree := h.buildPermissionTree(permissions, 0)
	response.Success(c, tree)
}

// GetPermission 获取权限详情
func (h *PermissionCRUDHandler) GetPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("权限ID格式错误"))
		return
	}

	permission, err := h.permissionRepo.Get(c.Request.Context(), uint(id), nil)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取权限详情失败", "id", id, "error", err)
		response.Error(c, err)
		return
	}

	response.Success(c, permission)
}

// CreatePermission 创建权限
func (h *PermissionCRUDHandler) CreatePermission(c *gin.Context) {
	var req dto.PermissionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	permission := &model.Permission{
		ParentID: req.ParentID,
		Name:     req.Name,
		Key:      req.Key,
		Type:     req.Type,
	}

	if err := h.permissionRepo.Create(c.Request.Context(), permission); err != nil {
		logger.ErrorContext(c.Request.Context(), "创建权限失败", "permission", permission, "error", err)
		response.Error(c, err)
		return
	}

	logger.InfoContext(c.Request.Context(), "创建权限成功", "permission_id", permission.ID, "key", permission.Key)
	response.Success(c, permission)
}

// UpdatePermission 更新权限
func (h *PermissionCRUDHandler) UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("权限ID格式错误"))
		return
	}

	var req dto.PermissionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg(err.Error()))
		return
	}

	permission, err := h.permissionRepo.Get(c.Request.Context(), uint(id), nil)
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "获取权限失败", "id", id, "error", err)
		response.Error(c, err)
		return
	}

	permission.ParentID = req.ParentID
	permission.Name = req.Name
	permission.Key = req.Key
	permission.Type = req.Type

	if err := h.permissionRepo.Update(c.Request.Context(), permission); err != nil {
		logger.ErrorContext(c.Request.Context(), "更新权限失败", "id", id, "error", err)
		response.Error(c, err)
		return
	}

	logger.InfoContext(c.Request.Context(), "更新权限成功", "permission_id", permission.ID, "key", permission.Key)
	response.Success(c, permission)
}

// DeletePermission 删除权限
func (h *PermissionCRUDHandler) DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("权限ID格式错误"))
		return
	}

	// 检查是否有子权限
	hasChildren, err := h.permissionRepo.HasChildren(c.Request.Context(), uint(id))
	if err != nil {
		logger.ErrorContext(c.Request.Context(), "检查子权限失败", "id", id, "error", err)
		response.Error(c, err)
		return
	}

	if hasChildren {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("该权限下还有子权限，无法删除"))
		return
	}

	if err := h.permissionRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		logger.ErrorContext(c.Request.Context(), "删除权限失败", "id", id, "error", err)
		response.Error(c, err)
		return
	}

	logger.InfoContext(c.Request.Context(), "删除权限成功", "permission_id", id)
	response.Success(c, struct{}{})
}

// buildPermissionTree 构建权限树
func (h *PermissionCRUDHandler) buildPermissionTree(permissions []model.Permission, parentID uint) []dto.PermissionTreeNode {
	var tree []dto.PermissionTreeNode

	for _, permission := range permissions {
		if permission.ParentID == parentID {
			node := dto.PermissionTreeNode{
				ID:       permission.ID,
				ParentID: permission.ParentID,
				Name:     permission.Name,
				Key:      permission.Key,
				Type:     permission.Type,
				Children: h.buildPermissionTree(permissions, permission.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}
