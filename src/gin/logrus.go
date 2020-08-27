package gin

import (
    "github.com/gin-gonic/gin"
    "github.com/mix-go/logrus"
    "time"
)

// LogrusWithFormatter instance a Logger middleware with the specified log format function.
func LogrusWithFormatter(logger *logrus.Logger, f gin.LogFormatter) gin.HandlerFunc {
    return func(c *gin.Context) {
        logger.Info(format(c, f))
    }
}

// LoggerWithConfig instance a Logger middleware with config.
func format(c *gin.Context, f gin.LogFormatter) string {
    // Start timer
    start := time.Now()
    path := c.Request.URL.Path
    raw := c.Request.URL.RawQuery

    // Process request
    c.Next()

    param := gin.LogFormatterParams{
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
