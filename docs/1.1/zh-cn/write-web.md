# 编写一个 Web 页面

帮助你快速搭建项目骨架，并指导你如何使用该骨架的细节。

## 安装

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
   ▸ Web (contains the websocket)
     gRPC
 ~~~

## 编写一个 Web 页面

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
├── public
├── routes
├── runtime
└── templates
~~~

`main.go` 文件：

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

我们可以在这里自定义命令，[查看更多](zh-cn/mix-xcli.md)

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
System      Name:      darwin
Go          Version:   1.13.4
Listen      Addr:      :8080
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
System      Name:      darwin
Go          Version:   1.13.4
Listen      Addr:      :8080
time=2020-09-16 20:24:41.515 level=info msg=Server start file=web.go:58
~~~

浏览器测试:

- 我们使用现成的工具测试：http://www.easyswoole.com/wstool.html

![](https://git.kancloud.cn/repos/onanying/mixgo1/raw/4f39803852f2155004172761ae95633f4666a3e5/images/websocket_test.png?access-token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2MTc5OTU4NjEsImlhdCI6MTYxNzk1MjY2MSwicmVwb3NpdG9yeSI6Im9uYW55aW5nXC9taXhnbzEiLCJ1c2VyIjp7InVzZXJuYW1lIjoib25hbnlpbmciLCJuYW1lIjoiXHU2NGI4XHU0ZWUzXHU3ODAxXHU3Njg0XHU0ZTYxXHU0ZTBiXHU0ZWJhIiwiZW1haWwiOiJjb2Rlci5saXVAcXEuY29tIiwidG9rZW4iOiIxODk5ZjEwODIzZWYwMmUxNzQ1MTgzMjk4YjhjNzFkMyIsImF1dGhvcml6ZSI6eyJwdWxsIjp0cnVlLCJwdXNoIjp0cnVlLCJhZG1pbiI6dHJ1ZX19fQ.sia533vNsLqbVetvwvttysxgWdRbbz0vfmN6jBNdl2g)

## 如何使用 DI 容器中的 Logger、Database、Redis 等组件

项目中要使用的公共组件，都定义在 `di` 目录，框架默认生成了一些常用的组件，用户也可以定义自己的组件，[查看更多](zh-cn/mix-xdi.md)

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
