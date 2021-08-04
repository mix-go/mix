> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XCLI

命令行交互指挥官

CLI Interactive Commander

## Overview

一个命令行交互与指挥管理工具，它可以让单个 CLI 可以执行多种功能，同时它还包括命令行参数获取、全局 panic 捕获与处理、程序后台执行等命令行开发常用功能。

## Installation

```
go get github.com/mix-go/xcli
```

## Quick start

```go
package main

import (
    "github.com/mix-go/xcli"
    "github.com/mix-go/xcli/flag"
)

func main() {
    xcli.SetName("app").SetVersion("0.0.0-alpha")
    cmd := &xcli.Command{
        Name:  "hello",
        Short: "Echo demo",
        RunF: func() {
            name := flag.Match("n", "name").String("default")
            // do something
        },
    }
    opt := &xcli.Option{
        Names: []string{"n", "name"},
        Usage: "Your name",
    }
    cmd.AddOption(opt)
    xcli.AddCommand(cmd).Run()
}
```

编译后，查看整个命令行程序的帮助

```
$ ./go_build_main_go 
Usage: ./go_build_main_go [OPTIONS] COMMAND [opt...]

Commands:
  hello         Echo demo

Global Options:
  -h, --help    Print usage
  -v, --version Print version information

Run './go_build_main_go COMMAND --help' for more information on a command.

Developed with Mix Go framework. (openmix.org/mix-go)
```

查看命令行程序的版本信息

```
$ ./go_build_main_go -v
app 0.0.0-alpha
```

查看 `hello` 命令的帮助

```
$ ./go_build_main_go hello --help
Usage: ./go_build_main_go hello [opt...]

Command Options:
  -n, --name    Your name

Developed with Mix Go framework. (openmix.org/mix-go)
```

执行 `hello` 命令

```
$ ./go_build_main_go hello 
```

## Flag 参数获取

> 该 flag 比 golang 自带的更加好用，不需要 Parse 操作

参数规则 (部分UNIX风格+GNU风格)

```
/examples/app home -d -rf --debug -v vvv --page 23 -s=test --name=john arg0
```
- 命令：
    - 第一个参数，可以为空：`home`
- 选项：
    - 短选项：一个中杠，如 `-d`、`-rf`
    - 长选项：二个中杠，如：`--debug`
- 选项值：
    - 无值：`-d`、`-rf`、 `--debug`
    - 有值(空格)：`-v vvv`、`--page 23`
    - 有值(等号)：`-s=test`、`--name=john`
- 参数：
    - 没有定义 `-` 的参数：`arg0`

获取选项，可以获取 `String`、`Bool`、`Int64`、`Float64` 多种类型，也可以指定默认值。

```
name := flag.Match("n", "name").String("Xiao Ming")
```

获取第一个参数

```
arg0 := flag.Arguments().First().String()
```

获取全部参数

```
for k, v := range flag.Arguments().Values() {
    // do something
}
```

## Daemon 后台执行

将命令行程序变为后台执行，该方法只可在 Main 协程中使用。

```
process.Daemon()
```

我们可以通过配合 `flag` 获取参数，实现通过某几个参数控制程序后台执行。

```go
if flag.Match("d", "daemon").Bool() {
    process.Daemon()
}
```

上面就实现了一个当命令行参数中带有 `-d/--daemon` 参数时，程序就在后台执行。

## Handle panic 错误处理

使用中间件处理异常，也可以单独对某个命令配置中间件

```go
h := func(next func()) {
    defer func() {
        if err := recover(); err != nil {
            // handle panic
        }
    }()
    next()
}
cmd := &xcli.Command{
    Name:  "hello",
    Short: "Echo demo",
    RunF: func() {
        // do something
    },
}
xcli.Use(h).AddCommand(cmd).Run()
```

## Application

我们在编写代码时，可能会要用到 App 中的一些信息。

```
// 获取基础路径(二进制所在目录路径)
xcli.App().BasePath

// App名称
xcli.App().Name

// App版本号
xcli.App().Version

// 是否开启debug
xcli.App().Debug
```

## Singleton 单命令

当我们的 CLI 只有一个命令时，只需要配置一下 `Singleton`：

- 也可以使用 `RunI` 执行一个结构体

~~~go
cmd := &xcli.Command{
    Name:  "hello",
    Short: "Echo demo",
    RunF: func() {
        // do something
    },
    Singleton: true,
}
~~~

命令的 Options 将会在 `-h/--help` 中打印

~~~
$ ./go_build_main_go 
Usage: ./go_build_main_go [OPTIONS] COMMAND [opt...]

Command Options:
  -n, --name    Your name

Global Options:
  -h, --help    Print usage
  -v, --version Print version information

Run './go_build_main_go --help' for more information on a command.

Developed with Mix Go framework. (openmix.org/mix-go)
~~~

## Default 默认执行

当我们的 CLI 有 CUI 时，需要实现点击后默认启动 UI 界面，只需要配置一下 `Default`：

~~~go
cmd := &xcli.Command{
    Name:  "hello",
    Short: "Echo demo",
    RunF: func() {
        // do something
    },
    Default: true,
}
~~~

## License

Apache License Version 2.0, http://www.apache.org/licenses/
