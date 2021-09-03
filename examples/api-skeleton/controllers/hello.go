package controllers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type HelloController struct {
}

func (t *HelloController) Index(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  http.StatusOK,
        "message": "hello, world!",
    })
}
