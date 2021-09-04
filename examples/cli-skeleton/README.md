## CLI development skeleton

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
   ▸ CLI
     API
     Web (contains the websocket)
     gRPC
 ~~~

## 编写一个 CLI 程序

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
├── logs
└── main.go
~~~

`main.go` 文件：

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

首先我们需要安装 [mix-go/xwp](https://github.com/mix-go/xwp)，因为这是一个独立库没有包含在骨架中：

```
go get github.com/mix-go/xwp
```

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
