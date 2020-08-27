package gin

import "github.com/gin-gonic/gin"

type IRouter interface {
    gin.IRouter
    Run() error
}

func New(definitions ...func(router *gin.Engine)) IRouter {
    e := gin.New()
    for _, d := range definitions {
        d(e)
    }
    return e
}
