package gin

import "github.com/gin-gonic/gin"

type RouteDefinition func(router *gin.Engine)

func New(definitions ...RouteDefinition) *gin.Engine {
    e := gin.New()
    for _, d := range definitions {
        d(e)
    }
    return e
}
