> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix CLI

一个快速创建 go 项目的脚手架

A scaffold to quickly create a go project

### Installation

- 安装

```
go get -u github.com/mix-go/mixcli
```

### Help

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

### New 

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

## License

Apache License Version 2.0, http://www.apache.org/licenses/
