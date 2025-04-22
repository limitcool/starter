package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/services"
	"github.com/spf13/cast"
)

// APIController API控制器
type APIController struct {
	apiService     *services.APIService
	menuAPIService *services.MenuAPIService
}

// NewAPIController 创建API控制器
func NewAPIController(apiService *services.APIService, menuAPIService *services.MenuAPIService) *APIController {
	controller := &APIController{
		apiService:     apiService,
		menuAPIService: menuAPIService,
	}

	// 将控制器添加到全局变量
	Controllers.APIController = controller

	return controller
}

// GetAPIs 获取API列表
func (c *APIController) GetAPIs(ctx *gin.Context) {
	apis, err := c.apiService.GetAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, apis)
}

// GetAPI 获取API详情
func (c *APIController) GetAPI(ctx *gin.Context) {
	id, err := cast.ToUint64E(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	api, err := c.apiService.GetByID(ctx.Request.Context(), uint(id))
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

	if err := c.apiService.Create(ctx.Request.Context(), &api); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, api)
}

// UpdateAPI 更新API
func (c *APIController) UpdateAPI(ctx *gin.Context) {
	id, err := cast.ToUint64E(ctx.Param("id"))
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
	if err := c.apiService.Update(ctx.Request.Context(), &api); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, api)
}

// DeleteAPI 删除API
func (c *APIController) DeleteAPI(ctx *gin.Context) {
	id, err := cast.ToUint64E(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	if err := c.apiService.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success[any](ctx, nil)
}

// AssignAPIsToMenu 为菜单分配API
func (c *APIController) AssignAPIsToMenu(ctx *gin.Context) {
	menuID, err := cast.ToUint64E(ctx.Param("id"))
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

	if err := c.menuAPIService.AssignAPIsToMenu(ctx.Request.Context(), uint(menuID), req.APIIDs); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success[any](ctx, nil)
}

// GetMenuAPIs 获取菜单关联的API
func (c *APIController) GetMenuAPIs(ctx *gin.Context) {
	menuID, err := cast.ToUint64E(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	apis, err := c.menuAPIService.GetMenuAPIs(ctx.Request.Context(), uint(menuID))
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, apis)
}

// SyncMenuAPIPermissions 同步菜单API权限
func (c *APIController) SyncMenuAPIPermissions(ctx *gin.Context) {
	if err := c.menuAPIService.SyncMenuAPIPermissions(ctx.Request.Context()); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success[any](ctx, nil)
}
