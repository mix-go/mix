> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix Console

命令行控制台程序开发框架

Command line console program development framework

> 该库还有 php 版本：https://github.com/mix-php/console

## Contents

- [Overview](#overview)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Flag](#flag)
- [Event](#event)
- [Daemon](#daemon)
- [Catch panic](#catch-panic)
- [Application](#application)
- [Bean](#bean)
- [Logger](#logger)

## Overview

Mix Console 不仅仅只是一个命令行骨架，它还包括命令行参数获取、依赖注入、事件驱动、全局 panic 捕获，错误信息接入日志或者第三方、程序后台执行等各种命令行开发常用功能，是一个完整的命令行程序开发框架。

## Installation

- 安装

```
go get -u github.com/mix-go/console
```

## Quick start

下面这些复杂的配置，在 [mix](https://github.com/mix-go/mix) 生成的骨架中都已经配置好了，并且做了合理的目录规划。

```go
package main

import (
    "github.com/mix-go/bean"
    "github.com/mix-go/console"
    "github.com/mix-go/event"
    "github.com/mix-go/logrus"
)

type HelloCommand struct {
}

func (t *HelloCommand) Main() {
    // do something
}

func main() {
    definition := console.ApplicationDefinition{
        AppName:    "app",
        AppVersion: "0.0.0-alpha",
        AppDebug:   true,
        // 该字段为程序依赖配置，内部的 eventDispatcher, error 是两个核心依赖，是必须配置的
        Beans: []bean.BeanDefinition{
            bean.BeanDefinition{
                Name:            "eventDispatcher",
                Reflect:         bean.NewReflect(event.NewDispatcher),
                Scope:           bean.SINGLETON,
                ConstructorArgs: bean.ConstructorArgs{},
            },
            bean.BeanDefinition{
                Name:            "error",
                Reflect:         bean.NewReflect(console.NewError),
                Scope:           bean.SINGLETON,
                ConstructorArgs: bean.ConstructorArgs{bean.NewReference("logger")},
                Fields: bean.Fields{
                    "Dispatcher": bean.NewReference("eventDispatcher"),
                },
            },
            bean.BeanDefinition{
                Name:    "logger",
                Reflect: bean.NewReflect(logrus.NewLogger),
                Scope:   bean.SINGLETON,
            },
        },
        // 该字段配置了程序有多少个命令，并且包含的参数，所有在程序中需要使用的参数都必须在这里定义
        Commands: []console.CommandDefinition{
            console.CommandDefinition{
                Name:  "hello",
                Usage: "\tEcho demo",
                Options: []console.OptionDefinition{
                    {
                        Names: []string{"n", "name"},
                        Usage: "Your name",
                    },
                    {
                        Names: []string{"say"},
                        Usage: "\tSay ...",
                    },
                },
                Command: &HelloCommand{},
                /* Singleton: true, // 如果这个程序只有一个命令，就开启这个配置 */
            },
        },
    }
    // 后面两个参数指定了必须配置的两个核心依赖的名称
    console.NewApplication(definition, "eventDispatcher", "error").Run()
}
```

编译后，查看整个命令行程序的帮助

```
$ ./go_build_main_go 
Usage: ./go_build_main_go [OPTIONS] COMMAND [opt...]

Global Options:
  -h, --help    Print usage
  -v, --version Print version information

Commands:
  hello         Echo demo


Run './go_build_main_go COMMAND --help' for more information on a command.

Developed with Mix Go framework. (openmix.org/mix-go)
```

查看命令行程序的版本信息

```
$ ./go_build_main_go -v
app 0.0.0-alpha, framework 1.0.9
```

查看 `hello` 命令的帮助

```
$ ./go_build_main_go hello --help
Usage: ./go_build_main_go hello [opt...]

Command Options:
  -n, --name    Your name
  --say         Say ...

Developed with Mix Go framework. (openmix.org/mix-go)
```

执行 `hello` 命令

```
$ ./go_build_main_go hello 
```

## Flag 

> 该 flag 比 golang 自带的更加好用，不需要 Parse 操作

获取命令行参数，可以获取 `String`、`Bool`、`Int64`、`Float64` 多种类型，也可以指定默认值。

```
name := flag.Match("n", "name").String("Xiao Ming")
```

参数规则 (部分UNIX风格+GNU风格)

- 单字母参数只支持一个中杠，如 `-p`，多字母参数只支持二个中杠，如：`--option`
- 参数可以有值、也可以没有值，如：
    - 无值：`-p`、 `--option`
    - 有值(空格)：`-p value`、`--option value`
    - 有值(等号)：`-p=value`、`--option=value`

## Event 

控制台基于 [Mix Event](https://github.com/mix-go/event) 管理本身的核心事件调度，让用户可以自定义处理

- `console.CommandBeforeExecuteEvent` 当命令执行前会调度该事件，用户可监听该事件在命令执行前处理前置逻辑，比如：将程序后台执行
- `console.HandleErrorEvent` 当程序捕获到全局 panic 或者代码中手动捕获的错误信息时调度该事件，用户可监听该事件将错误信息打印到日志或者发送到 Sentry 等平台。 

## Daemon

将命令行程序变为后台执行，该方法只可在 Main 协程中使用。

```
process.Daemon()
```

我们可以通过配合 `console.CommandBeforeExecuteEvent` 和 `flag` 获取参数，实现通过某几个参数控制程序后台执行。

```go
package listeners

import (
    "github.com/mix-go/console"
    "github.com/mix-go/console/flag"
    "github.com/mix-go/console/process"
    "github.com/mix-go/event"
)

type CommandListener struct {
}

func (t *CommandListener) Events() []event.Event {
    return []event.Event{
        &console.CommandBeforeExecuteEvent{},
    }
}

func (t *CommandListener) Process(e event.Event) {
    switch e.(type) {
    case *console.CommandBeforeExecuteEvent:
        // 设置守护
        if flag.Match("d", "daemon").Bool() {
            process.Daemon()
        }
        break
    }
}
```

创建的监听器需要注册到 `eventDispatcher` 组件的构造参数中

```
ConstructorArgs: bean.ConstructorArgs{listeners.CommandListener{}},
```

上面就实现了一个当命令行参数中带有 `-d/--daemon` 参数时，程序就在后台执行，注意：这两个参数需要在 `console.CommandDefinition` 中配置才可使用。

## Catch panic

go 程序的 err 返回设计虽然用户手动处理了大部分的错误，但是总是会有一些运行时 panic 是忘记处理的，但是这个错误信息是默认直接输出在 `os.Stdout` 中，日志中无法看到，非常容易忽略并且 debug 困难，我们解决了这个问题，可以将全部错误捕获集中交给 `error` 组件处理。

- Main 主协程
   - `error` 当主协程中的抛出 panic 时，Application 会内部使用了 recover 将错误信息传递到 `error` 组件处理。
- 子协程：子协程里的 panic 只能在子协程中使用 recover 捕获。
   - `catch.Error` 可以手动在子协程通过 recover 将错误信息使用该方法传递给 `error` 组件处理，如：`catch.Error(err, true)` 第二个参数强制打印堆栈信息。
   - `catch.Call` 更加省事的方法就是在开启子协程的位置使用该方法，如：`go foo()` 修改为 `go catch.Call(foo)`

错误信息将会集中到 `error` 处理，该组件会调用 `logger` 组件打印到日志，logger 组件会将 panic 的错误堆栈信息打印到日志中，这样就不用怕忽略错误信息，并且 debug 将变得更加容易，该组件还会使用 `eventDispatcher` 组件调度 `console.HandleErrorEvent` 事件，用户可将错误信息自定接入 Sentry 等第三方平台。

## Application

我们在 `Command: &HelloCommand{}` 中编写代码时，经常会要调用 App 中的一些功能。

```
console.App()
```

APP 的一些属性

```
// 获取基础路径(二进制所在目录路径)
console.App().BasePath

// App名称
console.App().AppName

// App版本号
console.App().AppVersion

// 是否开启debug
console.App().AppDebug

// 依赖注入容器
console.App().Context
```

由于依赖注入容器使用非常频繁，于是我们可以这样快速获取组件

```
console.Get("logger").(*logrus.Logger)
```

上面语句实际上等于这个

```
console.App().Context.Get("logger").(*logrus.Logger)
```

## Bean

控制台基于 [Mix Bean](https://github.com/mix-go/bean) 管理程序内部库的全部依赖关系，实现 DI、IoC，上面提到的 `console.App().Context` 就是控制台通过该库创建，该库设计思想参考 java spring，使用非常灵活。

## Logger

控制台 `error` 组件创建时必须传一个实现 Logger 接口的参数：

```go
type Logger interface {
    ErrorStack(err interface{}, stack *[]byte)
}
```

这个接口是为了实现可以通过 Logger 打印 panic 的堆栈信息到日志中，[Mix Logrus](https://github.com/mix-go/logrus) 实现了这个接口，因此可以直接接入，如需要使用其他开源 Logger 就需要用户自行扩展实现上面的接口。

## License

Apache License Version 2.0, http://www.apache.org/licenses/
