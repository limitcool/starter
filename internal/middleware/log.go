package middleware

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/limitcool/lib"
)

// LoggerWithCharmbracelet creates a custom Gin middleware that uses charmbracelet/log for logging.
func LoggerWithCharmbracelet() gin.HandlerFunc {
	var logger *log.Logger
	lib.SetDebugMode(func() {
		log.Info("Gin Debug Mode")
		logger = log.NewWithOptions(os.Stdout,
			log.Options{
				Level:           log.DebugLevel,
				Prefix:          "🚀 Gin",
				ReportTimestamp: true,
				ReportCaller:    true,
			},
		)
	})
	logger = log.New(os.Stdout)

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		stopTime := time.Since(startTime)
		latency := stopTime.Milliseconds()
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()

		fields := []interface{}{
			"status", statusCode,
			"method", method,
			"path", path,
			"query", query,
			"ip", clientIP,
			"user_agent", userAgent,
			"referer", referer,
			"latency_ms", latency,
		}

		if statusCode >= 500 {
			logger.Error("Request failed", fields...)
		} else if statusCode >= 400 {
			logger.Warn("Request processed with warning", fields...)
		} else {
			logger.Info("Request processed", fields...)
		}
	}
}
