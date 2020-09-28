## Mix Gin

基于 Gin 的 Web 库, 扩展 logrus 支持

Gin based Web library, extend Logrus support

## Overview

与原版 `Gin` 有哪些不同？

- 扩展了通过 `func(router *gin.Engine)` 闭包初始化路由的方式，方便路由定义规划。
- 扩展了路由日志对第三方 `logger` 的支持，包括：`logrus` 等。

## Installation

- 安装

```
go get -u github.com/mix-go/gin
```

## Usage

通过闭包创建路由

```go
routeDefinitionCallback := func(router *gin.Engine) {
    router.GET("hello",
        middleware.CorsMiddleware(),
        func(ctx *gin.Context) {
            hello := controllers.HelloController{}
            hello.Index(ctx)
        },
    )

    router.POST("users/add",
        middleware.CorsMiddleware(),
        func(ctx *gin.Context) {
            hello := controllers.AddUserController{}
            hello.Index(ctx)
        },
    )

    router.POST("auth", func(ctx *gin.Context) {
        auth := controllers.AuthController{}
        auth.Index(ctx)
    })
}
router := gin.New(routeDefinitionCallback)
```

接入第三方 `logrus` 打印日志

```go
logger := logrus.NewLogger()
router.Use(gin.LoggerWithFormatter(logger, func(params gin.LogFormatterParams) string {
    return fmt.Sprintf("%s|%s|%d|%s",
        params.Method,
        params.Path,
        params.StatusCode,
        params.ClientIP,
    )
}))
```

可以接入实现以下接口的所有第三方日志组件

```go
type Logger interface {
    Info(args ...interface{})
}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
