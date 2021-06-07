> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

<p align="center">
    <br>
    <br>
    <img src="https://openmix.org/static/image/logo_go.png" width="120" alt="MixPHP">
    <br>
    <br>
</p>

<h1 align="center">Mix Go</h1>

Mix Go 是一个基于 Go 进行快速开发的完整系统，类似前端的 `Vue CLI`，提供：

- 通过 `mix-go/mixcli` 实现的交互式项目脚手架：
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

如果编译时报错，整理一下依赖

~~~
go mod tidy
~~~

## 技术交流

知乎：https://www.zhihu.com/people/onanying   
微博：http://weibo.com/onanying    
官方QQ群：[284806582](https://shang.qq.com/wpa/qunwpa?idkey=b3a8618d3977cda4fed2363a666b081a31d89e3d31ab164497f53b72cf49968a), [825122875](http://shang.qq.com/wpa/qunwpa?idkey=d2908b0c7095fc7ec63a2391fa4b39a8c5cb16952f6cfc3f2ce4c9726edeaf20)，敲门暗号：goer

## 视频教程

[![从 PHP 转 Go 的基础知识对比视频讲解](https://openstr.com/cover/41e9dc609cb8f9a4530fe8f7a37f1130.jpg?size=small&share=true)](https://openstr.com/watch/41e9dc609cb8f9a4530fe8f7a37f1130)

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

我们可以在这里自定义命令，[查看更多](https://github.com/mix-go/xcli)

- `RunI` 定义了 `hello` 命令执行的接口，也可以使用 `RunF` 设定一个匿名函数

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

### 编写一个 Worker Pool 队列消费

队列消费是高并发系统中最常用的异步处理模型，通常我们是编写一个 CLI 命令行程序在后台执行 Redis、RabbitMQ 等 MQ 的队列消费，并将处理结果落地到 mysql 等数据库中，由于这类需求的标准化比较容易，因此我们开发了 [mix-go/xwp](https://github.com/mix-go/xwp) 库来处理这类需求，基本上大部分异步处理类需求都可使用。

新建 `commands/workerpool.go` 文件：

- `Foo` 结构体负责任务的执行处理，任务数据会在 `Do` 方法中触发，只需将我们的业务逻辑写到该方法中即可
- `p := &xwp.WorkerPool` 创建了一个协程池
- 当程序接收到进程退出信号时，协程池 `p.Stop()` 能平滑控制所有的 Worker 在执行完队列里全部的任务后再退出，保证数据的完整性

~~~go
package commands

import (
    "context"
    "fmt"
    "github.com/mix-go/cli-skeleton/di"
    "github.com/mix-go/xwp"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"
)

type Foo struct {
}

func (t *Foo) Do(data interface{}) {
    defer func() {
        if err := recover(); err != nil {
            logger := di.Logrus()
            logger.Error(err)
        }
    }()
    
    // 类型断言 str := data.(string)

    // 执行业务处理
    // ...
    
    // 将处理结果落地到数据库
    // ...
}

type WorkerPoolDaemonCommand struct {
}

func (t *WorkerPoolDaemonCommand) Main() {
    redis := di.GoRedis()
    jobQueue := make(chan interface{}, 200)
    p := &xwp.WorkerPool{
        JobQueue:       jobQueue,
        MaxWorkers:     1000,
        InitWorkers:    100,
        MaxIdleWorkers: 100,
        RunI:           &Foo{},
    }

    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-ch
        p.Stop()
    }()

    go func() {
        for {
            res, err := redis.BRPop(context.Background(), 3*time.Second, "foo").Result()
            if err != nil {
                if strings.Contains(err.Error(), "redis: nil") {
                    continue
                }
                fmt.Println(fmt.Sprintf("Redis Error: %s", err))
                p.Stop()
                return
            }
            // brPop命令最后一个键才是值
            jobQueue <- res[1]
        }
    }()

    p.Run() // 阻塞等待
}
~~~

接下来只需要把这个命令通过 `xcli.AddCommand` 注册到 CLI 中即可。

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
Listen      Addr:      :8080
System      Name:      darwin
Go          Version:   1.13.4
Framework   Version:   1.0.9
time=2020-09-16 20:24:41.515 level=info msg=Server start file=api.go:58
~~~

## 编写一个 Web 服务

首先我们使用 `mixcli` 命令创建一个项目骨架：

~~~
$ mixcli new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
    CLI
    API
  ▸ Web (contains the websocket)
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
├── main.go
├── middleware
├── public
├── routes
├── runtime
└── templates
~~~

`mian.go` 文件：

- `xcli.AddCommand` 方法传入的 `commands.Commands` 定义了全部的命令

~~~go
package main

import (
	"github.com/mix-go/web-skeleton/commands"
	_ "github.com/mix-go/web-skeleton/configor"
	_ "github.com/mix-go/web-skeleton/di"
	_ "github.com/mix-go/web-skeleton/dotenv"
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

~~~go
package commands

import (
	"github.com/mix-go/xcli"
)

var Commands = []*xcli.Command{
	{
		Name:  "web",
		Short: "\tStart the web server",
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
		RunI: &WebCommand{},
	},
}
~~~

`commands/web.go` 文件：

业务代码写在 `WebCommand` 结构体的 `main` 方法中，生成的代码中已经包含了：

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
	"github.com/mix-go/dotenv"
	"github.com/mix-go/web-skeleton/di"
	"github.com/mix-go/web-skeleton/routes"
	"github.com/mix-go/xcli"
	"github.com/mix-go/xcli/flag"
	"github.com/mix-go/xcli/process"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type WebCommand struct {
}

func (t *WebCommand) Main() {
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

	// templates
	router.LoadHTMLGlob(fmt.Sprintf("%s/../templates/*", xcli.App().BasePath))

	// static file
	router.Static("/static", fmt.Sprintf("%s/../public/static", xcli.App().BasePath))
	router.StaticFile("/favicon.ico", fmt.Sprintf("%s/../public/favicon.ico", xcli.App().BasePath))

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
	"github.com/mix-go/web-skeleton/controllers"
	"github.com/mix-go/web-skeleton/middleware"
)

func SetRoutes(router *gin.Engine) {
	router.Use(gin.Recovery()) // error handle

	router.GET("hello",
		func(ctx *gin.Context) {
			hello := controllers.HelloController{}
			hello.Index(ctx)
		},
	)

	router.Any("users/add",
		middleware.SessionMiddleware(),
		func(ctx *gin.Context) {
			user := controllers.UserController{}
			user.Add(ctx)
		},
	)

	router.Any("login", func(ctx *gin.Context) {
		login := controllers.LoginController{}
		login.Index(ctx)
	})

	router.GET("websocket",
		func(ctx *gin.Context) {
			ws := controllers.WebSocketController{}
			ws.Index(ctx)
		},
	)
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

命令行启动 `web` 服务器：

~~~
$ bin/go_build_main_go web
             ___         
 ______ ___  _ /__ ___ _____ ______ 
  / __ `__ \/ /\ \/ /__  __ `/  __ \
 / / / / / / / /\ \/ _  /_/ // /_/ /
/_/ /_/ /_/_/ /_/\_\  \__, / \____/ 
                     /____/


Server      Name:      mix-web
Listen      Addr:      :8080
System      Name:      darwin
Go          Version:   1.13.4
Framework   Version:   1.0.9
time=2020-09-16 20:24:41.515 level=info msg=Server start file=web.go:58
~~~

浏览器测试:

- 首先浏览器进入 http://127.0.0.1:8080/login 获取 session

![](https://git.kancloud.cn/repos/onanying/mixgo1/raw/4f39803852f2155004172761ae95633f4666a3e5/images/login.png?access-token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2MTc5OTU4NjEsImlhdCI6MTYxNzk1MjY2MSwicmVwb3NpdG9yeSI6Im9uYW55aW5nXC9taXhnbzEiLCJ1c2VyIjp7InVzZXJuYW1lIjoib25hbnlpbmciLCJuYW1lIjoiXHU2NGI4XHU0ZWUzXHU3ODAxXHU3Njg0XHU0ZTYxXHU0ZTBiXHU0ZWJhIiwiZW1haWwiOiJjb2Rlci5saXVAcXEuY29tIiwidG9rZW4iOiIxODk5ZjEwODIzZWYwMmUxNzQ1MTgzMjk4YjhjNzFkMyIsImF1dGhvcml6ZSI6eyJwdWxsIjp0cnVlLCJwdXNoIjp0cnVlLCJhZG1pbiI6dHJ1ZX19fQ.sia533vNsLqbVetvwvttysxgWdRbbz0vfmN6jBNdl2g)

- 提交表单后跳转到 http://127.0.0.1:8080/users/add 页面

![](https://git.kancloud.cn/repos/onanying/mixgo1/raw/4f39803852f2155004172761ae95633f4666a3e5/images/user-add.png?access-token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2MTc5OTU4NjEsImlhdCI6MTYxNzk1MjY2MSwicmVwb3NpdG9yeSI6Im9uYW55aW5nXC9taXhnbzEiLCJ1c2VyIjp7InVzZXJuYW1lIjoib25hbnlpbmciLCJuYW1lIjoiXHU2NGI4XHU0ZWUzXHU3ODAxXHU3Njg0XHU0ZTYxXHU0ZTBiXHU0ZWJhIiwiZW1haWwiOiJjb2Rlci5saXVAcXEuY29tIiwidG9rZW4iOiIxODk5ZjEwODIzZWYwMmUxNzQ1MTgzMjk4YjhjNzFkMyIsImF1dGhvcml6ZSI6eyJwdWxsIjp0cnVlLCJwdXNoIjp0cnVlLCJhZG1pbiI6dHJ1ZX19fQ.sia533vNsLqbVetvwvttysxgWdRbbz0vfmN6jBNdl2g)

### 编写一个 WebSocket 服务

WebSocket 是基于 http 协议完成握手的，因此我们编写代码时，也是和编写 Web 项目是差不多的，差别就是请求过来后，我们需要使用一个 WebSocket 的升级器，将请求升级为 WebSocket 连接，接下来就是针对连接的逻辑处理，从这个部分开始就和传统的 Socket 操作一致了。

`routes/main.go` 文件已经定义了一个 WebSocket 的路由：

~~~go
router.GET("websocket",
    func(ctx *gin.Context) {
        ws := controllers.WebSocketController{}
        ws.Index(ctx)
    },
)
~~~

`controllers/ws.go` 文件：

- 创建了一个 `upgrader` 的升级器，当请求过来时将会升级为 WebSocket 连接
- 定义了一个 `WebSocketSession` 的结构体负责管理连接的整个生命周期
- `session.Start()` 中启动了两个协程，分别处理消息的读和写
- 在消息读取的协程中，启动了 `WebSocketHandler` 结构体的 `Index` 方法来处理消息，在实际项目中我们可以根据不同的消息内容使用不同的结构体来处理，实现 Web 项目那种控制器的功能

~~~go
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mix-go/web-skeleton/di"
	"github.com/mix-go/xcli"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketController struct {
}

func (t *WebSocketController) Index(c *gin.Context) {
	logger := di.Logrus()
	if xcli.App().Debug {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	session := WebSocketSession{
		Conn:   conn,
		Header: c.Request.Header,
		Send:   make(chan []byte, 100),
	}
	session.Start()

	server := di.Server()
	server.RegisterOnShutdown(func() {
		session.Stop()
	})

	logger.Infof("Upgrade: %s", c.Request.UserAgent())
}

type WebSocketSession struct {
	Conn   *websocket.Conn
	Header http.Header
	Send   chan []byte
}

func (t *WebSocketSession) Start() {
	go func() {
		logger := di.Logrus()
		for {
			msgType, msg, err := t.Conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, 1001, 1006) {
					logger.Error(err)
				}
				t.Stop()
				return
			}
			if msgType != websocket.TextMessage {
				continue
			}

			handler := WebSocketHandler{
				Session: t,
			}
			handler.Index(msg)
		}
	}()
	go func() {
		logger := di.Logrus()
		for {
			msg, ok := <-t.Send
			if !ok {
				return
			}
			if err := t.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				logger.Error(err)
				t.Stop()
				return
			}
		}
	}()
}

func (t *WebSocketSession) Stop() {
	defer func() {
		if err := recover(); err != nil {
			logger := di.Logrus()
			logger.Error(err)
		}
	}()
	close(t.Send)
	_ = t.Conn.Close()
}

type WebSocketHandler struct {
	Session *WebSocketSession
}

func (t *WebSocketHandler) Index(msg []byte) {
	t.Session.Send <- []byte("hello, world!")
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

在命令行启动 `web` 服务器：

~~~
$ bin/go_build_main_go web
             ___         
 ______ ___  _ /__ ___ _____ ______ 
  / __ `__ \/ /\ \/ /__  __ `/  __ \
 / / / / / / / /\ \/ _  /_/ // /_/ /
/_/ /_/ /_/_/ /_/\_\  \__, / \____/ 
                     /____/


Server      Name:      mix-web
Listen      Addr:      :8080
System      Name:      darwin
Go          Version:   1.13.4
Framework   Version:   1.0.9
time=2020-09-16 20:24:41.515 level=info msg=Server start file=web.go:58
~~~

浏览器测试:

- 我们使用现成的工具测试：http://www.easyswoole.com/wstool.html

![](https://git.kancloud.cn/repos/onanying/mixgo1/raw/4f39803852f2155004172761ae95633f4666a3e5/images/websocket_test.png?access-token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2MTc5OTU4NjEsImlhdCI6MTYxNzk1MjY2MSwicmVwb3NpdG9yeSI6Im9uYW55aW5nXC9taXhnbzEiLCJ1c2VyIjp7InVzZXJuYW1lIjoib25hbnlpbmciLCJuYW1lIjoiXHU2NGI4XHU0ZWUzXHU3ODAxXHU3Njg0XHU0ZTYxXHU0ZTBiXHU0ZWJhIiwiZW1haWwiOiJjb2Rlci5saXVAcXEuY29tIiwidG9rZW4iOiIxODk5ZjEwODIzZWYwMmUxNzQ1MTgzMjk4YjhjNzFkMyIsImF1dGhvcml6ZSI6eyJwdWxsIjp0cnVlLCJwdXNoIjp0cnVlLCJhZG1pbiI6dHJ1ZX19fQ.sia533vNsLqbVetvwvttysxgWdRbbz0vfmN6jBNdl2g)

## 编写一个 gRPC 服务、客户端

首先我们使用 `mixcli` 命令创建一个项目骨架：

~~~
$ mixcli new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
    CLI
    API
    Web (contains the websocket)
  ▸ gRPC
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
├── main.go
├── protos
├── runtime
└── services
~~~

`mian.go` 文件：

- `xcli.AddCommand` 方法传入的 `commands.Commands` 定义了全部的命令

~~~go
package main

import (
	"github.com/mix-go/dotenv"
	"github.com/mix-go/grpc-skeleton/commands"
	_ "github.com/mix-go/grpc-skeleton/configor"
	_ "github.com/mix-go/grpc-skeleton/di"
	_ "github.com/mix-go/grpc-skeleton/dotenv"
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

- 定义了 `grpc:server`、`grpc:client` 两个子命令
- `RunI` 指定了命令执行的接口，也可以使用 `RunF` 设定一个匿名函数

```go
package commands

import (
	"github.com/mix-go/xcli"
)

var Commands = []*xcli.Command{
	{
		Name:  "grpc:server",
		Short: "gRPC server demo",
		Options: []*xcli.Option{
			{
				Names: []string{"d", "daemon"},
				Usage: "Run in the background",
			},
		},
		RunI: &GrpcServerCommand{},
	},
	{
		Name:  "grpc:client",
		Short: "gRPC client demo",
		RunI:  &GrpcClientCommand{},
	},
}
```

`protos/user.proto` 数据结构文件：

客户端与服务器端代码中都需要使用 `.proto` 生成的 go 代码，因为双方需要使用该数据结构通讯

- `.proto` 是 [gRPC](https://github.com/grpc/grpc) 通信的数据结构文件，采用 [protobuf](https://github.com/protocolbuffers/protobuf) 协议

~~~
syntax = "proto3";

package go.micro.grpc.user;
option go_package = "./;protos";

service User {
    rpc Add(AddRequest) returns (AddResponse) {}
}

message AddRequest {
    string Name = 1;
}

message AddResponse {
    int32 error_code = 1;
    string error_message = 2;
    int64 user_id = 3;
}
~~~

然后我们需要安装 gRPC 相关的编译程序：

- https://www.cnblogs.com/oolo/p/11840305.html#%E5%AE%89%E8%A3%85-grpc

接下来我们开始编译 `.proto` 文件：

- 编译成功后会在当前目录生成 `protos/user.pb.go` 文件

~~~
cd protos
protoc --go_out=plugins=grpc:. user.proto
~~~

`commands/server.go` 文件：

服务端代码写在 `GrpcServerCommand` 结构体的 `main` 方法中，生成的代码中已经包含了：

- 监听信号停止服务
- 可选的后台守护执行
- `pb.RegisterUserServer` 注册了一个默认服务，用户只需要扩展自己的服务即可

~~~go
package commands

import (
	"github.com/mix-go/dotenv"
	"github.com/mix-go/grpc-skeleton/di"
	pb "github.com/mix-go/grpc-skeleton/protos"
	"github.com/mix-go/grpc-skeleton/services"
	"github.com/mix-go/xcli/flag"
	"github.com/mix-go/xcli/process"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var netListener net.Listener

type GrpcServerCommand struct {
}

func (t *GrpcServerCommand) Main() {
	if flag.Match("d", "daemon").Bool() {
		process.Daemon()
	}

	addr := dotenv.Getenv("GIN_ADDR").String(":8080")
	logger := di.Logrus()

	// listen
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	netListener = listener

	// signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		logger.Info("Server shutdown")
		if err := listener.Close(); err != nil {
			panic(err)
		}
	}()

	// server
	s := grpc.NewServer()
	pb.RegisterUserServer(s, &services.UserService{})

	// run
	welcome()
	logger.Infof("Server run %s", addr)
	if err := s.Serve(listener); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}
}
~~~

`services/user.go` 文件：

服务端代码中注册的 `services.UserService{}` 服务代码如下：

只需要填充业务逻辑即可

```go
package services

import (
	"context"
	pb "github.com/mix-go/grpc-skeleton/protos"
)

type UserService struct {
}

func (t *UserService) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {
	// 执行数据库操作
	// ...

	resp := pb.AddResponse{
		ErrorCode:    0,
		ErrorMessage: "",
		UserId:       10001,
	}
	return &resp, nil
}
```

`commands/client.go` 文件：

客户端代码写在 `GrpcClientCommand` 结构体的 `main` 方法中，生成的代码中已经包含了：

- 通过环境配置获取服务端连接地址
- 设定了 `5s` 的执行超时时间

~~~go
package commands

import (
    "context"
    "fmt"
	"github.com/mix-go/dotenv"
	pb "github.com/mix-go/grpc-skeleton/protos"
    "google.golang.org/grpc"
    "time"
)

type GrpcClientCommand struct {
}

func (t *GrpcClientCommand) Main() {
    addr := dotenv.Getenv("GIN_ADDR").String(":8080")
    ctx, _ := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
    conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        panic(err)
    }
    defer func() {
        _ = conn.Close()
    }()
    cli := pb.NewUserClient(conn)
    req := pb.AddRequest{
        Name: "xiaoliu",
    }
    resp, err := cli.Add(ctx, &req)
    if err != nil {
        panic(err)
    }
    fmt.Println(fmt.Sprintf("Add User: %d", resp.UserId))
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

首先在命令行启动 `grpc:server` 服务器：

~~~
$ bin/go_build_main_go grpc:server
             ___         
 ______ ___  _ /__ ___ _____ ______ 
  / __ `__ \/ /\ \/ /__  __ `/  __ \
 / / / / / / / /\ \/ _  /_/ // /_/ /
/_/ /_/ /_/_/ /_/\_\  \__, / \____/ 
                     /____/


Server      Name:      mix-grpc
Listen      Addr:      :8080
System      Name:      darwin
Go          Version:   1.13.4
Framework   Version:   1.0.20
time=2020-11-09 15:08:17.544 level=info msg=Server run :8080 file=server.go:46
~~~

然后开启一个新的终端，执行下面的客户端命令与上面的服务器通信

~~~
$ bin/go_build_main_go grpc:client
Add User: 10001
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
