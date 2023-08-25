## Mix CLI

一个快速创建 Go 项目的脚手架

## Installation

```
go install -u github.com/mix-go/mixcli
```

## New Project

创建项目

- 可以生成 `cli`, `api`, `web`, `grpc` 多种项目代码，生成的代码开箱即用
- 可选择是否需要 `.env` 环境配置
- 可选择使用 `viper`, `configor` 加载 `.yml`, `.json`, `.toml` 等独立配置
- 可选择使用 `gorm`, `xorm` 的数据库
- 可选择使用 `zap`, `logrus` 的日志库

~~~
$ mixcli new hello
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select project type:
  ▸ CLI
    API
    Web (contains the websocket)
    gRPC
~~~
