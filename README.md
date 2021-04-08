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

- `RunI` 指定了命令执行的接口，也可以使用 `Run` 设定一个匿名函数

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

## 编写一个 gRPC 服务

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

我们可以在这里编辑命令，[查看更多](https://github.com/mix-go/xcli)

- 定义了 `grpc:server`、`grpc:client` 两个子命令
- `RunI` 指定了命令执行的接口，也可以使用 `Run` 设定一个匿名函数

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
option go_package = ".;protos";

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

var listener net.Listener

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
	listener = listener

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

`commands/client.go` 文件：

客户端代码写在 `GrpcClientCommand` 结构体的 `main` 方法中，生成的代码中已经包含了：

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
Add User: 1200
~~~

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
