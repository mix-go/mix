package controllers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type HelloController struct {
}

func (t *HelloController) Index(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl", gin.H{
        "title": "Hello, World!",
    })
}
