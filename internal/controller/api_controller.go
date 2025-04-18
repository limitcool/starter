package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/services"
)

// APIController API控制器
type APIController struct {
	apiService     *services.APIService
	menuAPIService *services.MenuAPIService
}

// NewAPIController 创建API控制器
func NewAPIController(apiService *services.APIService, menuAPIService *services.MenuAPIService) *APIController {
	return &APIController{
		apiService:     apiService,
		menuAPIService: menuAPIService,
	}
}

// GetAPIs 获取API列表
func (c *APIController) GetAPIs(ctx *gin.Context) {
	apis, err := c.apiService.GetAll()
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, apis)
}

// GetAPI 获取API详情
func (c *APIController) GetAPI(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	api, err := c.apiService.GetByID(uint(id))
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, api)
}

// CreateAPI 创建API
func (c *APIController) CreateAPI(ctx *gin.Context) {
	var api model.API
	if err := ctx.ShouldBindJSON(&api); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	if err := c.apiService.Create(&api); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, api)
}

// UpdateAPI 更新API
func (c *APIController) UpdateAPI(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	var api model.API
	if err := ctx.ShouldBindJSON(&api); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	api.ID = uint(id)
	if err := c.apiService.Update(&api); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, api)
}

// DeleteAPI 删除API
func (c *APIController) DeleteAPI(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	if err := c.apiService.Delete(uint(id)); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success[any](ctx, nil)
}

// AssignAPIsToMenu 为菜单分配API
func (c *APIController) AssignAPIsToMenu(ctx *gin.Context) {
	menuID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	var req struct {
		APIIDs []uint `json:"api_ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	if err := c.menuAPIService.AssignAPIsToMenu(uint(menuID), req.APIIDs); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success[any](ctx, nil)
}

// GetMenuAPIs 获取菜单关联的API
func (c *APIController) GetMenuAPIs(ctx *gin.Context) {
	menuID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	apis, err := c.menuAPIService.GetMenuAPIs(uint(menuID))
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, apis)
}

// SyncMenuAPIPermissions 同步菜单API权限
func (c *APIController) SyncMenuAPIPermissions(ctx *gin.Context) {
	if err := c.menuAPIService.SyncMenuAPIPermissions(); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success[any](ctx, nil)
}
