package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
}

func (t *UserController) Add(c *gin.Context) {
	// 网页
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "user_add.tmpl", gin.H{
			"title": "User add",
		})
		c.Abort()
		return
	}

	// 执行数据库操作
	// ...

	c.String(http.StatusInternalServerError, "<html><h1>%s</h1></html>", "Add ok!")
}
