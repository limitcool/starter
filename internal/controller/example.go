package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/internal/pkg/apiresponse"
	"github.com/limitcool/starter/internal/pkg/code"
)

// ExampleHandler 是演示错误处理的示例控制器
func ExampleHandler(c *gin.Context) {
	var user model.SysUser
	db := services.Instance().GetDB()

	// 1. 使用apiresponse.HandleError处理错误
	if err := db.First(&user, 9999).Error; err != nil {
		// 这会自动处理GORM的ErrRecordNotFound错误
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success(c, user)
}

// ExampleErrorHandler 演示使用gin的错误处理机制
func ExampleErrorHandler(c *gin.Context) {
	var user model.SysUser
	db := services.Instance().GetDB()

	// 2. 使用gin的Error方法
	if err := db.First(&user, 9999).Error; err != nil {
		// 这会被ErrorHandler中间件捕获并处理
		_ = c.Error(err)
		return
	}

	apiresponse.Success(c, user)
}

// ExampleCustomError 演示使用自定义错误
func ExampleCustomError(c *gin.Context) {
	id := c.Param("id")

	// 3. 使用自定义错误码
	if id == "" {
		// 直接使用自定义错误
		apiresponse.HandleError(c, code.NewErrCode(code.InvalidParams))
		return
	}

	// 或者通过gin的Error方法
	if id == "0" {
		_ = c.Error(code.NewErrMsg("ID不能为0"))
		return
	}

	apiresponse.Success(c, gin.H{"id": id})
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
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success(c, user)
}
