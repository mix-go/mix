package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
}

func (t *UserController) Add(c *gin.Context) {
	// 执行数据库操作
	// ...

    c.JSON(http.StatusOK, gin.H{
        "status":  http.StatusOK,
        "message": "ok",
    })
}
