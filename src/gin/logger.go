package gin

import (
	"github.com/gin-gonic/gin"
	"time"
)

// LogFormatter 日志格式化处理器
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams 日志格式化参数
type LogFormatterParams gin.LogFormatterParams

// Logger 日志处理器接口
type Logger interface {
	Info(args ...interface{})
}

// LoggerWithFormatter 配置日志格式
func LoggerWithFormatter(logger Logger, f LogFormatter) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info(format(c, f))
	}
}

// LoggerWithConfig instance a Logger middleware with config.
func format(c *gin.Context, f LogFormatter) string {
	// Start timer
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	// Process request
	c.Next()

	param := LogFormatterParams{
		Request: c.Request,
		Keys:    c.Keys,
	}

	// Stop timer
	param.TimeStamp = time.Now()
	param.Latency = param.TimeStamp.Sub(start)

	param.ClientIP = c.ClientIP()
	param.Method = c.Request.Method
	param.StatusCode = c.Writer.Status()
	param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

	param.BodySize = c.Writer.Size()

	if raw != "" {
		path = path + "?" + raw
	}

	param.Path = path

	return f(param)
}
