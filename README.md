> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

<p align="center">
    <br>
    <br>
    <img src="https://openmix.org/static/image/logo_go.png" width="120" alt="MixPHP">
    <br>
    <br>
</p>

<h1 align="center">Mix Go</h1>

MixGo 是一个 Go 快速开发标准工具包；内部模块高度解耦，整体代码基于多个独立的模块构建，即便用户不使用我们的 `mixcli` 脚手架快速生成代码，也可以使用这些独立模块。例如：你可以只使用 `xcli` 来构建你的命令行交互；可以使用 `xsql` 来调用数据库；可以使用 `xwp` 来处理 MQ 队列消费；所有的模块你可以像搭积木一样随意组合。

## 独立模块

核心模块全部可独立使用。

- [mix-go/mixcli](src/mixcli) 快速创建 Go 项目的脚手架，类似前端界的 Vue CLI。
- [mix-go/xcli](src/xcli) 命令行交互与指挥管理工具，同时它还包括命令行参数获取、中间件、程序守护等。
- [mix-go/xsql](src/xsql) 基于 database/sql 的轻量数据库，功能完备且支持任何数据库驱动。
- [mix-go/xdi](src/xdi) 处理对象依赖关系的 IoC、DI 库，可以实现统一管理依赖，全局对象管理，动态配置刷新等。
- [mix-go/xwp](src/xwp) 一个通用工作池、协程池，可动态扩容缩容。
- [mix-go/xutil](src/xutil) 一套让 Golang 保持甜美的工具。

## 开发文档

- `V1.1` https://openmix.org/mix-go/docs/1.1/
- `V1.0` https://www.kancloud.cn/onanying/mixgo1/content

## 快速开始

提供了现成的脚手架工具，快速创建项目，立即产出。

- [编写一个 CLI 程序](examples/cli-skeleton#readme)
  - [编写一个 Worker Pool 队列消费](examples/cli-skeleton#%E7%BC%96%E5%86%99%E4%B8%80%E4%B8%AA-worker-pool-%E9%98%9F%E5%88%97%E6%B6%88%E8%B4%B9)
- [编写一个 API 服务](examples/api-skeleton#readme)
- [编写一个 Web 服务](examples/web-skeleton#readme)
  - [编写一个 WebSocket 服务](examples/web-skeleton#%E7%BC%96%E5%86%99%E4%B8%80%E4%B8%AA-WebSocket-%E6%9C%8D%E5%8A%A1)
- [编写一个 gRPC 服务、客户端](examples/grpc-skeleton#readme)


```bash
go install github.com/mix-go/mixcli@latest
```


```bash
$ mixcli new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
  ▸ CLI
    API
    Web (contains the websocket)
    gRPC
```

如果编译时报错，整理一下依赖

~~~
go mod tidy
~~~

## 推荐阅读

- [MixGo 在 IDE Goland 中的如何使用](https://zhuanlan.zhihu.com/p/391857663)

## 技术交流

知乎：https://www.zhihu.com/people/onanying    
官方QQ群：[284806582](https://shang.qq.com/wpa/qunwpa?idkey=b3a8618d3977cda4fed2363a666b081a31d89e3d31ab164497f53b72cf49968a), [825122875](http://shang.qq.com/wpa/qunwpa?idkey=d2908b0c7095fc7ec63a2391fa4b39a8c5cb16952f6cfc3f2ce4c9726edeaf20) 敲门暗号：gopher

## PHP 框架

OpenMix 同时还有 PHP 生态的框架

- https://github.com/mix-php/mix
- https://gitee.com/mix-php/mix

## License

Apache License Version 2.0, http://www.apache.org/licenses/
