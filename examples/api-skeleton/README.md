## API development skeleton

帮助你快速搭建项目骨架，并指导你如何使用该骨架的细节。

## Installation

- Install

~~~
go get -u github.com/mix-go/mixcli
~~~

- New project

~~~
mixcli new hello
~~~

~~~
 Use the arrow keys to navigate: ↓ ↑ → ← 
 ? Select project type:
     CLI
   ▸ API
     Web (contains the websocket)
     gRPC
 ~~~

## 编写一个 API 服务

首先我们使用 `mixcli` 命令创建一个项目骨架：

~~~
$ mixcli new hello
~~~

生成骨架目录结构如下：

~~~
.
├── README.md
├── bin
├── commands
├── conf
├── config
├── controllers
├── di
├── go.mod
├── go.sum
├── main.go
├── middleware
├── routes
└── runtime
~~~

`main.go` 文件：

- `xcli.AddCommand` 方法传入的 `commands.Commands` 定义了全部的命令

~~~go
package main

import (
	"github.com/mix-go/api-skeleton/commands"
	_ "github.com/mix-go/api-skeleton/configor"
	_ "github.com/mix-go/api-skeleton/di"
	_ "github.com/mix-go/api-skeleton/dotenv"
	"github.com/mix-go/dotenv"
	"github.com/mix-go/xcli"
)

func main() {
	xcli.SetName("app").
		SetVersion("0.0.0-alpha").
		SetDebug(dotenv.Getenv("APP_DEBUG").Bool(false))
	xcli.AddCommand(commands.Commands...).Run()
}
~~~

`commands/main.go` 文件：

我们可以在这里自定义命令，[查看更多](https://github.com/mix-go/xcli)

- `RunI` 指定了命令执行的接口，也可以使用 `RunF` 设定一个匿名函数

```go
package commands

import (
	"github.com/mix-go/xcli"
)

var Commands = []*xcli.Command{
	{
		Name:  "api",
		Short: "\tStart the api server",
		Options: []*xcli.Option{
			{
				Names: []string{"a", "addr"},
				Usage: "\tListen to the specified address",
			},
			{
				Names: []string{"d", "daemon"},
				Usage: "\tRun in the background",
			},
		},
		RunI: &APICommand{},
	},
}
```

`commands/api.go` 文件：

业务代码写在 `APICommand` 结构体的 `main` 方法中，生成的代码中已经包含了：

- 监听信号停止服务
- 根据模式打印日志
- 可选的后台守护执行

基本上无需修改即可上线使用

~~~go
package commands

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mix-go/api-skeleton/di"
	"github.com/mix-go/api-skeleton/routes"
	"github.com/mix-go/dotenv"
	"github.com/mix-go/xcli/flag"
	"github.com/mix-go/xcli/process"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type APICommand struct {
}

func (t *APICommand) Main() {
	if flag.Match("d", "daemon").Bool() {
		process.Daemon()
	}

	logger := di.Logrus()
	server := di.Server()
	addr := dotenv.Getenv("GIN_ADDR").String(":8080")
	mode := dotenv.Getenv("GIN_MODE").String(gin.ReleaseMode)

	// server
	gin.SetMode(mode)
	router := gin.New()
	routes.SetRoutes(router)
	server.Addr = flag.Match("a", "addr").String(addr)
	server.Handler = router

	// signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		logger.Info("Server shutdown")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			logger.Errorf("Server shutdown error: %s", err)
		}
	}()

	// logger
	if mode != gin.ReleaseMode {
		handlerFunc := gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: func(params gin.LogFormatterParams) string {
				return fmt.Sprintf("%s|%s|%d|%s",
					params.Method,
					params.Path,
					params.StatusCode,
					params.ClientIP,
				)
			},
			Output: logger.Out,
		})
		router.Use(handlerFunc)
	}

	// run
	welcome()
	logger.Infof("Server start at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !strings.Contains(err.Error(), "http: Server closed") {
		panic(err)
	}
}
~~~

在 `routes/main.go` 文件中配置路由：

已经包含一些常用实例，只需要在这里新增路由即可开始开发

~~~go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mix-go/api-skeleton/controllers"
	"github.com/mix-go/api-skeleton/middleware"
)

func SetRoutes(router *gin.Engine) {
	router.Use(gin.Recovery()) // error handle

	router.GET("hello",
		middleware.CorsMiddleware(),
		func(ctx *gin.Context) {
			hello := controllers.HelloController{}
			hello.Index(ctx)
		},
	)

	router.POST("users/add",
		middleware.AuthMiddleware(),
		func(ctx *gin.Context) {
			hello := controllers.UserController{}
			hello.Add(ctx)
		},
	)

	router.POST("auth", func(ctx *gin.Context) {
		auth := controllers.AuthController{}
		auth.Index(ctx)
	})
}
~~~

接下来我们编译上面的程序：

- linux & macOS

~~~
go build -o bin/go_build_main_go main.go
~~~

- win

~~~
go build -o bin/go_build_main_go.exe main.go
~~~

启动服务器

~~~
$ bin/go_build_main_go api
             ___         
 ______ ___  _ /__ ___ _____ ______ 
  / __ `__ \/ /\ \/ /__  __ `/  __ \
 / / / / / / / /\ \/ _  /_/ // /_/ /
/_/ /_/ /_/_/ /_/\_\  \__, / \____/ 
                     /____/


Server      Name:      mix-api
System      Name:      darwin
Go          Version:   1.13.4
Listen      Addr:      :8080
time=2020-09-16 20:24:41.515 level=info msg=Server start file=api.go:58
~~~

## 如何使用 DI 容器中的 Logger、Database、Redis 等组件

项目中要使用的公共组件，都定义在 `di` 目录，框架默认生成了一些常用的组件，用户也可以定义自己的组件，[查看更多](https://github.com/mix-go/xdi)

- 可以在哪里使用

可以在代码的任意位置使用，但是为了可以使用到环境变量和自定义配置，通常我们在 `xcli.Command` 结构体定义的 `RunF`、`RunI` 中使用。

- 使用日志，比如：`logrus`、`zap`

```go
logger := di.Logrus()
logger.Info("test")
```

- 使用数据库，比如：`gorm`、`xorm`

```go
db := di.Gorm()
user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
result := db.Create(&user)
fmt.Println(result)
```

- 使用 Redis，比如：`go-redis`

```go
rdb := di.GoRedis()
val, err := rdb.Get(context.Background(), "key").Result()
if err != nil {
    panic(err)
}
fmt.Println("key", val)
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
