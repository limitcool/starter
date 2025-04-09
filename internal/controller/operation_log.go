package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/request"
	"github.com/limitcool/starter/internal/pkg/apiresponse"
	"github.com/limitcool/starter/internal/services"
	"github.com/spf13/cast"
)

// GetOperationLogs 获取操作日志列表
func GetOperationLogs(c *gin.Context) {
	// 构建查询参数
	var query request.OperationLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		apiresponse.ParamError(c, "无效的查询参数")
		return
	}

	// 如果分页参数没有设置，使用默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = request.DefaultPageSize
	}

	// 调用服务查询数据
	db := services.Instance().GetDB()
	logService := services.NewOperationLogService(db)
	result, err := logService.GetOperationLogs(&query)
	if err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success(c, result)
}

// DeleteOperationLog 删除操作日志
func DeleteOperationLog(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		apiresponse.ParamError(c, "无效的ID参数")
		return
	}

	db := services.Instance().GetDB()
	logService := services.NewOperationLogService(db)
	if err := logService.DeleteOperationLog(cast.ToUint(id)); err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success[any](c, nil)
}

// BatchDeleteOperationLogs 批量删除操作日志
func BatchDeleteOperationLogs(c *gin.Context) {
	var req request.OperationLogBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ParamError(c, "无效的请求参数")
		return
	}

	db := services.Instance().GetDB()
	logService := services.NewOperationLogService(db)
	if err := logService.BatchDeleteOperationLogs(req.IDs); err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success[any](c, nil)
}
