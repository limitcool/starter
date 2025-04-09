package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/pkg/apiresponse"
	"github.com/limitcool/starter/pkg/errors"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 获取错误码和消息
			errCode, errMsg := errors.ParseError(err)
			apiresponse.Fail(c, errCode, errMsg)

			c.Abort()
		}
	}
}
