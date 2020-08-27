package gin

import "github.com/gin-gonic/gin"

func New(definitions ...func(router *gin.Engine)) *gin.Engine {
    e := gin.New()
    for _, d := range definitions {
        d(e)
    }
    return e
}

const (
    // DebugMode indicates gin mode is debug.
    DebugMode = "debug"
    // ReleaseMode indicates gin mode is release.
    ReleaseMode = "release"
    // TestMode indicates gin mode is test.
    TestMode = "test"
)

func SetMode(value string) {
    gin.SetMode(value)
}
