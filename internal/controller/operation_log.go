package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/dto"
	"github.com/limitcool/starter/internal/pkg/util"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/pkg/apiresponse"
)

// GetOperationLogs 获取操作日志列表
// @Summary 获取操作日志列表
// @Description 分页获取操作日志列表，支持多条件查询
// @Tags 操作日志
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页大小" default(20)
// @Param sort_by query string false "排序字段" default(id)
// @Param sort_desc query bool false "是否降序" default(true)
// @Param user_type query string false "用户类型" Enums(sys_user, user)
// @Param username query string false "用户名"
// @Param module query string false "操作模块"
// @Param action query string false "操作类型"
// @Param ip query string false "IP地址"
// @Param start_time query string false "开始时间" format(datetime)
// @Param end_time query string false "结束时间" format(datetime)
// @Success 200 {object} apiresponse.Response
// @Router /api/v1/admin/operation-logs [get]
func GetOperationLogs(c *gin.Context) {
	// 构建查询参数
	var query dto.OperationLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		apiresponse.ParamError(c, "无效的查询参数")
		return
	}

	// 如果分页参数没有设置，使用默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = dto.DefaultPageSize
	}

	// 调用服务查询数据
	logService := services.NewOperationLogService(global.DB)
	result, err := logService.GetOperationLogs(&query)
	if err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success(c, result)
}

// DeleteOperationLog 删除操作日志
// @Summary 删除操作日志
// @Description 根据ID删除单条操作日志
// @Tags 操作日志
// @Accept json
// @Produce json
// @Param id path int true "操作日志ID"
// @Success 200 {object} apiresponse.Response
// @Router /api/v1/admin/operation-logs/{id} [delete]
func DeleteOperationLog(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		apiresponse.ParamError(c, "无效的ID参数")
		return
	}

	logService := services.NewOperationLogService(global.DB)
	if err := logService.DeleteOperationLog(util.ParseUint(id, 0)); err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success[any](c, nil)
}

// BatchDeleteOperationLogs 批量删除操作日志
// @Summary 批量删除操作日志
// @Description 根据ID数组批量删除操作日志
// @Tags 操作日志
// @Accept json
// @Produce json
// @Param ids body dto.OperationLogBatchDeleteDTO true "操作日志ID数组"
// @Success 200 {object} apiresponse.Response
// @Router /api/v1/admin/operation-logs/batch [delete]
func BatchDeleteOperationLogs(c *gin.Context) {
	var req dto.OperationLogBatchDeleteDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ParamError(c, "无效的请求参数")
		return
	}

	logService := services.NewOperationLogService(global.DB)
	if err := logService.BatchDeleteOperationLogs(req.IDs); err != nil {
		apiresponse.ServerError(c)
		return
	}

	apiresponse.Success[any](c, nil)
}
