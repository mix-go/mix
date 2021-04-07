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
  - 可选择是否需要 `.yal`, `.json`, `.toml` 等独立配置
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
$> mixcli new hello
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
$> mixcli new hello
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

- `commands.Commands` 定义了全部的命令

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
				Short: "Your name",
			},
			{
				Names: []string{"say"},
				Short: "\tSay ...",
			},
		},
		RunI: &HelloCommand{},
	},
}
```

`HelloCommand` 文件：

业务代码写在该结构体的 `main` 方法中

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

~~~
// linux & macOS
go build -o bin/go_build_main_go main.go

// win
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

## 编写一个 Web 服务

## 编写一个 WebSocket 服务

## 编写一个 gRPC 服务、客户端

## 编写一个 Worker Pool 队列消费

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
