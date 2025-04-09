package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/request"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/services"
)

// GetOperationLogs 获取操作日志列表
func GetOperationLogs(c *gin.Context) {
	// 解析查询参数
	var query request.OperationLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ParamError(c, "无效的查询参数")
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
	db := services.Instance().GetDB()
	opLogService := services.NewOperationLogService(db)
	result, err := opLogService.GetOperationLogs(&query)
	if err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, result)
}

// DeleteOperationLog 删除操作日志
func DeleteOperationLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的ID参数")
		return
	}

	db := services.Instance().GetDB()
	opLogService := services.NewOperationLogService(db)
	if err := opLogService.DeleteOperationLog(uint(id)); err != nil {
		response.ServerError(c)
		return
	}

	response.Success[any](c, nil)
}

// ClearOperationLogs 清空操作日志
func ClearOperationLogs(c *gin.Context) {
	var req request.OperationLogBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "无效的请求参数")
		return
	}

	db := services.Instance().GetDB()
	opLogService := services.NewOperationLogService(db)
	if err := opLogService.BatchDeleteOperationLogs(req.IDs); err != nil {
		response.ServerError(c)
		return
	}

	response.Success[any](c, nil)
}
