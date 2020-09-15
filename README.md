> OpenMix 出品：https://openmix.org

<br>

<p align="center">
    <img src="https://openmix.org/static/image/logo_go.png" width="120" alt="MixPHP">
</p>

<p align="center">高性能 • 轻量级 • 命令行</p>

## Mix Go 是什么

Mix Go 是混合型高性能 Go 框架，该框架可以开发 console, api, web 等各种项目，引入了依赖注入、控制反转、事件驱动等高级特征，得益于 go 生态更好的跨平台、静态执行的优势，该框架更适合系统核心模块、对稳定性要求高、计算量比较大的项目。

## 与 Mix PHP 的关系

该框架与 [Mix PHP](https://github.com/mix-php/mix) 设计哲学几乎完全一致，Mix PHP 的用户可以非常容易的切换到 Mix Go 进行开发，达到学一会二的效果，OpenMix 可能是现在唯一一个打造跨语言框架的开源机构。

## 与其他 Go 框架的差别

- 骨架代码全部基于 bean, event 依赖注入、控制反转、事件驱动库构建，同时内置了 Go 生态各个领域最流行的库，包括 gin, gorm, logrus 等，我们已经将这些离散的库整合为一体，可以相互关联使用。

- 框架内置了 gin 作为服务器，gin 严格来讲并不是框架，而是一个 server 库，只提供了服务器相关的功能，请求处理，中间件，视图渲染等。

- 提供了 console, api, web 多种骨架生成工具，同时骨架代码中包含非常丰富的范例，开箱即用。

- 框架非常轻量灵活，依赖库均可独立使用，严格来讲除了 `mix-go/console` 命令行开发组件，其他全部为选装。

- 采用高度灵活的开发方式，框架只提供底层库，而与具体功能相关的代码都在骨架代码中实现，用户能更加细粒度的修改每一处细节。

- 由于骨架和核心类库都由 Mix 自己打造，拥有和 Mix PHP 同样的设计哲学，PHP 的用户可以很容易上手开发。

## 微服务

由于 Gin 与 go-micro 是兼容的，因此可以非常方便的扩展为微服务。

## 框架定位

当我们开发 Mix PHP 时发现框架的设计哲学可以复制到 Go 生态，于是我们着手实现让更多的 PHP 中级程序员也可使用 Go 打造高并发系统，让 Mix 的用户能学一会二，实现跨语言无差别开发。

## 开发文档

- https://openmix.org/mix-go/doc
- https://www.kancloud.cn/onanying/mixgo1/content

## 快速开始

- 安装开发工具

~~~
go get -u github.com/mix-go/mix@master
~~~

- 创建应用骨架

~~~
// console
mix new --name=hello
~~~

~~~
// api
mix api --name=hello
~~~

- 编译到骨架的 `bin` 目录

~~~
cd hello
go build -o bin/go_build_main_go main.go
~~~

- 执行

~~~
cd bin
./go_build_main_go
~~~

- `api` 测试

~~~
$> ./go_build_main_go api
             ___         
 ______ ___  _ /__ ___ _____ ______ 
  / __ `__ \/ /\ \/ /__  __ `/  __ \
 / / / / / / / /\ \/ _  /_/ // /_/ /
/_/ /_/ /_/_/ /_/\_\  \__, / \____/ 
                     /____/


Server      Name:     mix-api
System      Name:     darwin
Go          Version:  1.13.4
Framework   Version:  1.0.5
Listen      Addr:     :8080
time=2020-08-28 18:54:31 level=info msg=Server start file=api.go:58
~~~

访问测试 (新开一个终端)：

```
$> curl http://127.0.0.1:8080/hello
{"message":"hello, world!","status":200}
```

## 技术交流

知乎：https://www.zhihu.com/people/onanying   
微博：http://weibo.com/onanying    
官方QQ群：[284806582](https://shang.qq.com/wpa/qunwpa?idkey=b3a8618d3977cda4fed2363a666b081a31d89e3d31ab164497f53b72cf49968a), [825122875](http://shang.qq.com/wpa/qunwpa?idkey=d2908b0c7095fc7ec63a2391fa4b39a8c5cb16952f6cfc3f2ce4c9726edeaf20)，敲门暗号：goer

## License

Apache License Version 2.0, http://www.apache.org/licenses/
