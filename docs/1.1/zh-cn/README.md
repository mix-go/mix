> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

<p align="center">
    <br>
    <br>
    <img src="https://openmix.org/static/image/logo_go.png" width="120" alt="MixPHP">
    <br>
    <br>
</p>

<h1 align="center">Mix Go</h1>

## 简介

MixGo 是一个 Go 快速开发标准工具包；内部模块高度解耦，整体代码基于多个独立的模块构建，即便用户不使用我们的 mixcli 脚手架快速生成代码，也可以使用这些独立模块。例如：你可以只使用 xcli 来构建你的命令行交互；可以使用 xdi 来管理全局对象的依赖；可以使用 xwp 来处理 MQ 队列消费；所有的模块你可以像搭积木一样随意组合。

## 请帮忙 Star 一下

- https://github.com/mix-go/mix
- https://gitee.com/mix-go/mix

## 独立模块

核心模块全部可独立使用。

- [mix-go/mixcli](zh-cn/mix-mixcli) 快速创建 Go 项目的脚手架，类似前端界的 Vue CLI
- [mix-go/xcli](zh-cn/mix-xcli) 命令行交互与指挥管理工具，同时它还包括命令行参数获取、中间件、程序守护等。
- [mix-go/xdi](zh-cn/mix-xdi) 处理对象依赖关系的 IoC、DI 库，可以实现统一管理依赖，全局对象管理，动态配置刷新等。
- [mix-go/xwp](zh-cn/mix-xwp) 一个通用工作池、协程池，可动态扩容缩容。
- [mix-go/xfmt](zh-cn/mix-xfmt) 可以打印结构体嵌套指针地址内部数据的格式化库
- [mix-go/varwatch](zh-cn/mix-varwatch) 监视配置结构体变量的数据变化并执行一些任务
- [mix-go/dotenv](zh-cn/mix-dotenv) 具有类型转换功能的 DotEnv 环境配置库

## PHP 框架

OpenMix 同时还有 PHP 生态的框架

- https://github.com/mix-php/mix
- https://gitee.com/mix-php/mix

## 旧版文档

- `V1.0` https://www.kancloud.cn/onanying/mixgo1/content
