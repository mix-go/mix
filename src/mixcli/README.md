> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix CLI

一个快速创建 go 项目的脚手架

A scaffold to quickly create a go project

## Installation

```
go get github.com/mix-go/mixcli
```

## Help

查看命令帮助

~~~
$ mixcli
Usage: mixcli [OPTIONS] COMMAND [opt...]

Commands:
  new           Create a project

Global Options:
  -h, --help    Print usage
  -v, --version Print version information


Run 'mixcli COMMAND --help' for more information on a command.

Developed with Mix Go framework. (openmix.org/mix-go)
~~~

## New 

创建项目

- 可以生成 `cli`, `api`, `web`, `grpc` 多种项目代码，生成的代码开箱即用
- 可选择是否需要 `.env` 环境配置
- 可选择是否需要 `.yml`, `.json`, `.toml` 等独立配置
- 可选择使用 `gorm`, `xorm` 的数据库
- 可选择使用 `logrus`, `zap` 的日志库

~~~
$ mixcli new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
  ▸ CLI
    API
    Web (contains the websocket)
    gRPC
~~~

## License

Apache License Version 2.0, http://www.apache.org/licenses/
