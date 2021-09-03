package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func CorsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Headers", "Origin, Accept, Keep-Alive, User-Agent, Cache-Control, Content-Type, X-Requested-With, Authorization")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
        if c.Request.Method == "OPTIONS" {
            c.String(http.StatusOK, "")
            c.Abort()
            return
        }

        c.Next()
    }
}
