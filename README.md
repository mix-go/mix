> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

<p align="center">
    <img src="https://openmix.org/static/image/logo_go.png" width="120" alt="MixPHP">
</p>

## Mix Go

Mix Go 是一个基于 Go 进行快速开发的完整系统，类似前端的 `Vue CLI`，提供：

- 通过 `mix-go/xstart` 实现的交互式的项目脚手架：
  - 可以生成 `cli`, `api`, `web`, `grpc` 多种项目代码
  - 生成的代码开箱即用
  - 可选择是否需要 `.env` 环境配置
  - 可选择是否需要 `YAML`, `JSON`, `TOML` 等独立配置
  - 可选择使用 `gorm`, `xorm` 的数据库
  - 可选择使用 `logrus`, `zap` 的日志库
- 通过 `mix-go/xcli` 实现的命令行原型开发。
- 基于 `mix-go/xdi` 的 DI, IoC 容器。

## 官方库

- https://github.com/mix-go/xstart
- https://github.com/mix-go/xcli
- https://github.com/mix-go/xdi
- https://github.com/mix-go/xwp
- https://github.com/mix-go/xfmt
- https://github.com/mix-go/dotenv

## 第三方库

- https://github.com/gin-gonic/gin
- https://gorm.io/gorm
- https://github.com/go-redis/redis/v8
- https://github.com/jinzhu/configor
- https://github.com/sirupsen/logrus
- https://github.com/lestrrat-go/file-rotatelogs
- https://github.com/go-session/session
- https://github.com/go-session/redis
- https://github.com/gorilla/websocket
- https://github.com/dgrijalva/jwt-go
- https://github.com/go-playground/validator/v10
- https://github.com/golang/protobuf

## 快速开始

```
go get -u github.com/mix-go/xstart
```

~~~
$ xstart new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
  ▸ CLI
    API
    Web (contains the websocket)
    gRPC
~~~

## 开发文档

- https://openmix.org/mix-go/doc
- https://www.kancloud.cn/onanying/mixgo1/content

## 技术交流

知乎：https://www.zhihu.com/people/onanying   
微博：http://weibo.com/onanying    
官方QQ群：[284806582](https://shang.qq.com/wpa/qunwpa?idkey=b3a8618d3977cda4fed2363a666b081a31d89e3d31ab164497f53b72cf49968a), [825122875](http://shang.qq.com/wpa/qunwpa?idkey=d2908b0c7095fc7ec63a2391fa4b39a8c5cb16952f6cfc3f2ce4c9726edeaf20)，敲门暗号：goer

## License

Apache License Version 2.0, http://www.apache.org/licenses/
