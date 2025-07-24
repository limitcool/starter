package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/errspec"
	"github.com/limitcool/starter/internal/pkg/errorx"
)

// Result 定义API响应结构
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorHandler 处理错误的中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 检查是否是应用错误
			var appErr *errorx.AppError
			if !errors.As(err, appErr) {
				// 对于非应用错误，返回通用错误
				appErr = errspec.ErrUnknown.New(c.Request.Context())
			}

			c.JSON(appErr.HttpStatus(), Result{
				Code:    appErr.Code(),
				Message: appErr.Error(),
				Data:    nil,
			})
			c.Abort()
		}
	}
}

// UserLoginHandler 模拟用户登录处理
func UserLoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	ctx := c.Request.Context()

	// 参数验证
	if username == "" || password == "" {
		c.Error(errspec.ErrUserNameOrPasswordEmpty.New(ctx))
		return
	}

	// 模拟用户查询
	if username != "admin" {
		c.Error(errspec.ErrUserNotFound.New(ctx))
		return
	}

	// 模拟密码验证
	if password != "123456" {
		c.Error(errspec.ErrPassword.New(ctx))
		return
	}

	// 登录成功
	c.JSON(http.StatusOK, Result{
		Code:    errspec.Success.Code(),
		Message: "登录成功",
		Data:    map[string]string{"username": username},
	})
}

func main() {

	ctx := context.TODO()

	// 演示直接使用错误
	fmt.Println("错误示例:", errspec.ErrUserNotFound.New(ctx))

	// 演示使用GetError
	unknownErr := errspec.ErrInternal.New(ctx)
	fmt.Println("通过错误码获取错误:", unknownErr.Error())

	// 演示添加额外信息
	customErr := errspec.ErrNotFound.New(ctx).WithMessage("用户ID为123的用户")
	fmt.Println("自定义错误消息:", customErr.Error())

	// 设置Gin路由
	r := gin.Default()
	r.Use(ErrorHandler())

	r.POST("/login", UserLoginHandler)

	fmt.Println("Server running at :8080")
	r.Run(":8080")
}
