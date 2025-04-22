package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/services"
	"github.com/spf13/cast"
)

func NewOperationLogController(opLogService *services.OperationLogService) *OperationLogController {
	return &OperationLogController{
		opLogService: opLogService,
	}
}

type OperationLogController struct {
	opLogService *services.OperationLogService
}

// GetOperationLogs 获取操作日志列表
func (olc *OperationLogController) GetOperationLogs(c *gin.Context) {
	// 解析查询参数
	var query v1.OperationLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 默认分页参数
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	// 获取操作日志
	result, err := olc.opLogService.GetOperationLogs(c.Request.Context(), &query)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// DeleteOperationLog 删除操作日志
func (olc *OperationLogController) DeleteOperationLog(c *gin.Context) {
	id, err := cast.ToUint64E(c.Param("id"))
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	if err := olc.opLogService.DeleteOperationLog(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// ClearOperationLogs 批量删除操作日志
func (olc *OperationLogController) ClearOperationLogs(c *gin.Context) {
	var req v1.OperationLogBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 检查ID列表是否为空
	if len(req.IDs) == 0 {
		response.Error(c, errorx.ErrInvalidParams.WithMsg("ID列表不能为空"))
		return
	}

	// 批量删除操作日志
	if err := olc.opLogService.BatchDeleteOperationLogs(c.Request.Context(), req.IDs); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}
