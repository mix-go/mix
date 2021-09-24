## gRPC development skeleton

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
     API
     Web (contains the websocket)
   ▸ gRPC
 ~~~

## 编写一个 gRPC 服务、客户端

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
├── di
├── go.mod
├── go.sum
├── main.go
├── protos
├── runtime
└── services
~~~

`main.go` 文件：

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
System      Name:      darwin
Go          Version:   1.13.4
Listen      Addr:      :8080
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

- 使用日志，比如：[zap](https://github.com/uber-go/zap)、[logrus](https://github.com/Sirupsen/logrus)

```go
logger := di.Zap()
logger.Info("test")
```

- 使用数据库，比如：[gorm](https://gorm.io/)、[xorm](https://xorm.io/)

```go
db := di.Gorm()
user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
result := db.Create(&user)
fmt.Println(result)
```

- 使用 Redis，比如：[go-redis](https://redis.uptrace.dev/)

```go
rdb := di.GoRedis()
val, err := rdb.Get(context.Background(), "key").Result()
if err != nil {
panic(err)
}
fmt.Println("key", val)
```

## 部署

线上部署时，不需要部署源码到服务器，只需要部署编译好的二进制、配置文件等

```
├── bin
├── conf
├── runtime
├── shell
└── .env
```

修改 `shell/server.sh` 脚本中的绝对路径和参数

```
file=/project/bin/program
cmd=grpc:server
```

启动管理

```
sh shell/server.sh start
sh shell/server.sh stop
sh shell/server.sh restart
```

gRPC 通常都是内部使用，使用内网 `SLB` 代理到服务器IP或者直接使用 IP:PORT 调用

## License

Apache License Version 2.0, http://www.apache.org/licenses/
