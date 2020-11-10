package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/mix-go/logrus"
	"time"
)

type LogFormatter func(params LogFormatterParams) string
type LogFormatterParams gin.LogFormatterParams

type Logger interface {
	Info(args ...interface{})
}

func LoggerWithFormatter(logger Logger, f LogFormatter) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info(format(c, f))
	}
}

// Deprecated: 使用 LoggerWithFormatter 替代
func LogrusWithFormatter(logger *logrus.Logger, f LogFormatter) gin.HandlerFunc {
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
