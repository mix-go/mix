## Mix Console

命令行控制台程序开发框架

Command line console program development framework

> 该库还有 php 版本：https://github.com/mix-php/console

## Overview

Mix Console 不仅仅只是一个命令行骨架，它还包括命令行参数获取、依赖注入、事件驱动、全局 panic 捕获，错误信息接入日志或者第三方、程序后台执行等各种命令行开发常用功能，是一个完整的命令行程序开发框架。

## Installation

- 安装

```
go get -u github.com/mix-go/console
```

## Quick start

```
definition := console.ApplicationDefinition{
    AppName:    "app",
    AppVersion: "0.0.0-alpha",
    AppDebug:   true,
    // 该字段为程序依赖配置，内部的 event, error 是两个核心依赖，是必须配置的
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
            Command: &commands.HelloCommand{},
            /* Singleton: true, // 如果这个程序只有一个命令，就开启这个配置 */
        },
    },
}
// 后面两个参数指定了必须配置的两个核心依赖的名称
console.NewApplication(definition, "eventDispatcher", "error").Run()
```

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

```
$ ./go_build_main_go -v
app 0.0.0-alpha, framework 1.0.9
```

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

```
$ ./go_build_main_go hello --help
Usage: ./go_build_main_go hello [opt...]

Command Options:
  -n, --name    Your name
  --say         Say ...

Developed with Mix Go framework. (openmix.org/mix-go)
```

## Flag 

获取命令行参数，可以获取 `String`、`Bool`、`Int64`、`Float64` 多种类型，也可以指定默认值。

```
name := flag.Match("n", "name").String("Xiao Ming")
```

## Event 

控制台基于 [Mix Event](https://github.com/mix-go/event) 管理本身的核心事件调度，让用户可以自定义处理

- `console.CommandBeforeExecuteEvent` 当命令执行前会调度该事件，用户可监听该事件在命令执行前处理前置逻辑，比如：将程序后台执行
- `console.HandleErrorEvent` 当程序捕获到全局 panic 或者代码中手动捕获的错误信息时调度该事件，用户可监听该事件将错误信息打印到日志或者发送到 Sentry 等平台。 

## Catch error



## License

Apache License Version 2.0, http://www.apache.org/licenses/
