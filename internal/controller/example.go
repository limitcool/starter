package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/services"
)

// ExampleHandler 示例处理函数
func ExampleHandler(c *gin.Context) {
	// 1. 使用响应工具
	// 直接返回成功响应，携带数据
	response.Success(c, gin.H{
		"message": "这是一个示例API",
		"status":  "success",
	})

	// 2. 更多响应示例
	// response.SuccessWithMsg(c, "自定义消息", data)
	// response.Fail(c, code.ErrorUnknown, "发生错误")
	// response.ParamError(c, "参数错误")
	// response.ServerError(c, "服务器错误")
	// response.NotFound(c, "资源不存在")
	// response.Unauthorized(c, "未授权")
	// response.Forbidden(c, "禁止访问")
}

// ExampleErrorHandler 演示错误处理
func ExampleErrorHandler(c *gin.Context) {
	// 模拟产生一个错误
	err := errorx.NewErrMsg("这是一个示例错误")
	// 使用统一错误处理
	response.HandleError(c, err)
}

// ExampleCustomError 演示使用自定义错误
func ExampleCustomError(c *gin.Context) {
	id := c.Param("id")

	// 3. 使用自定义错误码
	if id == "" {
		// 直接使用自定义错误
		response.HandleError(c, errorx.NewErrCode(errorx.InvalidParams))
		return
	}

	// 或者通过gin的Error方法
	if id == "0" {
		_ = c.Error(errorx.NewErrMsg("ID不能为0"))
		return
	}

	response.Success(c, gin.H{"id": id})
}

// ExampleDBOperationError 演示数据库操作错误处理
func ExampleDBOperationError(c *gin.Context) {
	// 4. 处理数据库操作错误
	user := model.SysUser{
		Username: "admin", // 假设这是唯一键，已存在
	}
	db := services.Instance().GetDB()

	// 插入会失败，因为用户名已存在(唯一约束)
	if err := db.Create(&user).Error; err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, user)
}
