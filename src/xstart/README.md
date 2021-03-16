> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XStart

A scaffold to quickly create a go project

一个快速创建 go 项目的脚手架

### Installation

- 安装

```
go get -u github.com/mix-go/xstart
```

### Help

查看命令帮助

~~~
$ xstart
Usage: xstart [OPTIONS] COMMAND [opt...]

Commands:
  new           Create a project

Global Options:
  -h, --help    Print usage
  -v, --version Print version information


Run 'xstart COMMAND --help' for more information on a command.

Developed with Mix Go framework. (openmix.org/mix-go)
~~~

### New 

创建项目

~~~
$ xstart new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
  ▸ CLI
    API
    Web (contains the websocket)
    gRPC
~~~

## License

Apache License Version 2.0, http://www.apache.org/licenses/
