package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/pkg/errorx"
)

// Response 定义API响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
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
			if errorx.IsAppErr(err) {
				appErr = err.(*errorx.AppError)
			} else {
				// 对于非应用错误，返回通用错误
				appErr = errorx.ErrUnknown
			}

			c.JSON(appErr.GetHttpStatus(), Response{
				Code: appErr.GetErrorCode(),
				Msg:  appErr.GetErrorMsg(),
				Data: nil,
			})
			c.Abort()
		}
	}
}

// UserLoginHandler 模拟用户登录处理
func UserLoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 参数验证
	if username == "" || password == "" {
		c.Error(errorx.ErrInvalidParams.WithMsg("用户名或密码不能为空"))
		return
	}

	// 模拟用户查询
	if username != "admin" {
		c.Error(errorx.ErrUserNotFound)
		return
	}

	// 模拟密码验证
	if password != "123456" {
		c.Error(errorx.ErrUserPasswordError)
		return
	}

	// 登录成功
	c.JSON(http.StatusOK, Response{
		Code: errorx.SuccessCode,
		Msg:  "登录成功",
		Data: map[string]string{"username": username},
	})
}

func main() {
	// 演示直接使用错误
	fmt.Println("错误示例:", errorx.ErrUserNotFound.Error())

	// 演示使用GetError
	unknownErr := errorx.GetError(10001)
	fmt.Println("通过错误码获取错误:", unknownErr.Error())

	// 演示添加额外信息
	customErr := errorx.ErrNotFound.WithMsg("用户ID为123的用户")
	fmt.Println("自定义错误消息:", customErr.Error())

	// 设置Gin路由
	r := gin.Default()
	r.Use(ErrorHandler())

	r.POST("/login", UserLoginHandler)

	// 错误码查询接口
	r.GET("/error/:code", func(c *gin.Context) {
		code := c.Param("code")
		var codeInt int
		fmt.Sscanf(code, "%d", &codeInt)

		err := errorx.GetError(codeInt)
		c.JSON(http.StatusOK, Response{
			Code: err.GetErrorCode(),
			Msg:  err.GetErrorMsg(),
			Data: nil,
		})
	})

	fmt.Println("Server running at :8080")
	r.Run(":8080")
}
