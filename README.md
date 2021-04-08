> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

<p align="center">
    <br>
    <img src="https://openmix.org/static/image/logo_go.png" width="120" alt="MixPHP">
    <br>
</p>

## Mix Go

Mix Go 是一个基于 Go 进行快速开发的完整系统，类似前端的 `Vue CLI`，提供：

- 通过 `mix-go/mixcli` 实现的交互式的项目脚手架：
  - 可以生成 `cli`, `api`, `web`, `grpc` 多种项目代码
  - 生成的代码开箱即用
  - 可选择是否需要 `.env` 环境配置
  - 可选择是否需要 `.yml`, `.json`, `.toml` 等独立配置
  - 可选择使用 `gorm`, `xorm` 的数据库
  - 可选择使用 `logrus`, `zap` 的日志库
- 通过 `mix-go/xcli` 实现的命令行原型开发。
- 基于 `mix-go/xdi` 的 DI, IoC 容器。

## 快速开始

安装

```
go get github.com/mix-go/mixcli
```

创建项目

~~~
$ mixcli new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
  ▸ CLI
    API
    Web (contains the websocket)
    gRPC
~~~

## 技术交流

知乎：https://www.zhihu.com/people/onanying   
微博：http://weibo.com/onanying    
官方QQ群：[284806582](https://shang.qq.com/wpa/qunwpa?idkey=b3a8618d3977cda4fed2363a666b081a31d89e3d31ab164497f53b72cf49968a), [825122875](http://shang.qq.com/wpa/qunwpa?idkey=d2908b0c7095fc7ec63a2391fa4b39a8c5cb16952f6cfc3f2ce4c9726edeaf20)，敲门暗号：goer

## 编写一个 CLI 程序

首先我们使用 `mixcli` 命令创建一个项目骨架：

~~~
$ mixcli new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
  ▸ CLI
    API
    Web (contains the websocket)
    gRPC
~~~

生成骨架目录结构如下：

~~~
.
├── README.md
├── bin
├── commands
├── conf
├── configor
├── di
├── dotenv
├── go.mod
├── go.sum
├── logs
└── main.go
~~~

`mian.go` 文件：

- `xcli.AddCommand` 方法传入的 `commands.Commands` 定义了全部的命令

~~~go
package main

import (
	"github.com/mix-go/cli-skeleton/commands"
	_ "github.com/mix-go/cli-skeleton/configor"
	_ "github.com/mix-go/cli-skeleton/di"
	_ "github.com/mix-go/cli-skeleton/dotenv"
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

我们可以在这里编辑命令，[查看更多](https://github.com/mix-go/xcli)

- `RunI` 定义了 `hello` 命令执行的接口，也可以使用 `Run` 设定一个匿名函数

```go
package commands

import (
	"github.com/mix-go/xcli"
)

var Commands = []*xcli.Command{
	{
		Name:  "hello",
		Short: "\tEcho demo",
		Options: []*xcli.Option{
			{
				Names: []string{"n", "name"},
				Usage: "Your name",
			},
			{
				Names: []string{"say"},
				Usage: "\tSay ...",
			},
		},
		RunI: &HelloCommand{},
	},
}
```

`commands/hello.go` 文件：

业务代码写在 `HelloCommand` 结构体的 `main` 方法中

- 代码中可以使用 `flag` 获取命令行参数，[查看更多](https://github.com/mix-go/xcli#flag)

```go
package commands

import (
	"fmt"
	"github.com/mix-go/xcli/flag"
)

type HelloCommand struct {
}

func (t *HelloCommand) Main() {
	name := flag.Match("n", "name").String("OpenMix")
	say := flag.Match("say").String("Hello, World!")
	fmt.Printf("%s: %s\n", name, say)
}
```

接下来我们编译上面的程序：

- linux & macOS

~~~
go build -o bin/go_build_main_go main.go
~~~

- win

~~~
go build -o bin/go_build_main_go.exe main.go
~~~

查看全部命令的帮助信息：

~~~
$ cd bin
$ ./go_build_main_go 
Usage: ./go_build_main_go [OPTIONS] COMMAND [opt...]

Global Options:
  -h, --help    Print usage
  -v, --version Print version information

Commands:
  hello         Echo demo

Run './go_build_main_go COMMAND --help' for more information on a command.

Developed with Mix Go framework. (openmix.org/mix-go)
~~~

查看上面编写的 hello 命令的帮助信息：

~~~
$ ./go_build_main_go hello --help
Usage: ./go_build_main_go hello [opt...]

Command Options:
  -n, --name    Your name
  --say         Say ...

Developed with Mix Go framework. (openmix.org/mix-go)
~~~

执行 `hello` 命令，并传入两个参数：

~~~
$ ./go_build_main_go hello --name=liujian --say=hello
liujian: hello
~~~

## 编写一个 API 服务

首先我们使用 `mixcli` 命令创建一个项目骨架：

~~~
$ mixcli new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
    CLI
  ▸ API
    Web (contains the websocket)
    gRPC
~~~

生成骨架目录结构如下：

~~~
.
├── README.md
├── bin
├── commands
├── conf
├── configor
├── controllers
├── di
├── dotenv
├── go.mod
├── go.sum
├── logs
├── main.go
├── middleware
├── routes
└── runtime
~~~


`mian.go` 文件：

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

我们可以在这里编辑命令，[查看更多](https://github.com/mix-go/xcli)

- `RunI` 定义了 `hello` 命令执行的接口，也可以使用 `Run` 设定一个匿名函数

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
Listen      Addr:      :8080
System      Name:      darwin
Go          Version:   1.13.4
Framework   Version:   1.0.9
time=2020-09-16 20:24:41.515 level=info msg=Server start file=api.go:58
~~~

## 编写一个 Web 服务

## 编写一个 WebSocket 服务

## 编写一个 gRPC 服务、客户端

## 编写一个 Worker Pool 队列消费

## 如何使用 DI 容器中的 Logger、Database、Redis 等组件

项目中要使用的公共组件，都定义在 `di` 目录，框架默认生成了一些常用的组件，用户也可以定义自己的组件，[查看更多](https://github.com/mix-go/xdi)

- 可以在哪里使用

可以在代码的任意位置使用，但是为了可以使用到环境变量和自定义配置，通常我们在 `xcli.Command` 结构体定义的 `Run`、`RunI` 中使用。

- 使用日志，比如：`logrus`

```go
logger := di.Logrus()
logger.Info("test")
```

- 使用数据库，比如：`gorm`

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

## 依赖

官方库

- https://github.com/mix-go/mixcli
- https://github.com/mix-go/xcli
- https://github.com/mix-go/xdi
- https://github.com/mix-go/xwp
- https://github.com/mix-go/xfmt
- https://github.com/mix-go/dotenv

第三方库

- https://github.com/gin-gonic/gin
- https://gorm.io
- https://github.com/go-redis/redis
- https://github.com/jinzhu/configor
- https://github.com/uber-go/zap
- https://github.com/sirupsen/logrus
- https://github.com/natefinch/lumberjack
- https://github.com/lestrrat-go/file-rotatelogs
- https://github.com/go-session/session
- https://github.com/go-session/redis
- https://github.com/dgrijalva/jwt-go
- https://github.com/gorilla/websocket
- https://github.com/golang/grpc
- https://github.com/golang/protobuf

## License

Apache License Version 2.0, http://www.apache.org/licenses/
