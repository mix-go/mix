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
$ mixcli new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
  ▸ CLI
    API
    Web (contains the websocket)
    gRPC
~~~

## 编写一个 CLI 程序

## 编写一个 API 服务

## 编写一个 Web 服务

## 编写一个 WebSocket 服务

## 编写一个 gRPC 服务、客户端

## 编写一个 Worker Pool 队列消费

## 技术交流

知乎：https://www.zhihu.com/people/onanying   
微博：http://weibo.com/onanying    
官方QQ群：[284806582](https://shang.qq.com/wpa/qunwpa?idkey=b3a8618d3977cda4fed2363a666b081a31d89e3d31ab164497f53b72cf49968a), [825122875](http://shang.qq.com/wpa/qunwpa?idkey=d2908b0c7095fc7ec63a2391fa4b39a8c5cb16952f6cfc3f2ce4c9726edeaf20)，敲门暗号：goer

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
