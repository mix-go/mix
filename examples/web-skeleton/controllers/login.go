package controllers

import (
    "context"
    "github.com/gin-gonic/gin"
    "github.com/mix-go/web-skeleton/di"
    "net/http"
)

type LoginController struct {
}

func (t *LoginController) Index(c *gin.Context) {
    // 网页
    if c.Request.Method == http.MethodGet {
        c.HTML(http.StatusOK, "login.tmpl", gin.H{
            "title": "Login",
        })
        c.Abort()
        return
    }

    // 检查用户登录代码
    // ...

    // session
    session := di.Session()
    store, err := session.Start(context.Background(), c.Writer, c.Request)
    if err != nil {
        panic(err)
    }
    store.Set("userinfo", gin.H{
        "user_id": 10008,
    })
    if err := store.Save(); err != nil {
        panic(err)
    }

    // 跳转到登录成功页
    c.Redirect(http.StatusMovedPermanently, "/users/add")
}
