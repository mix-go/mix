> Produced by OpenMix: [https://openmix.org](https://openmix.org/mix-go)

<p align="center">
    <br>
    <br>
    <img src="https://openmix.org/static/image/logo_go.png" width="120" alt="MixPHP">
    <br>
    <br>
</p>

<h1 align="center">Mix Go</h1>

English | [中文](README_CN.md)

MixGo is a Go rapid development standard toolkit; the internal modules are highly decoupled, and the overall code is built on multiple independent modules. Even if users do not use our `mixcli` scaffolding to quickly generate code, they can also use these independent modules. For example: you can use `xcli` alone to build your command-line interaction; use `xsql` to call the database; use `xwp` to handle MQ queue consumption; you can freely combine all modules like building blocks.

## Independent Modules

All core modules can be used independently.

- [mix-go/mixcli](src/mixcli) Scaffold to quickly create Go projects, similar to Vue CLI in the frontend field.
- [mix-go/xcli](src/xcli) Command-line interaction and command management tool, also includes command-line parameter acquisition, middleware, program daemon, etc.
- [mix-go/xsql](src/xsql) Lightweight database based on database/sql, fully functional and supports any database driver.
- [mix-go/xdi](src/xdi) IoC, DI library for handling object dependencies, can implement unified dependency management, global object management, dynamic configuration refresh, etc.
- [mix-go/xwp](src/xwp) A universal work pool, coroutine pool, can dynamically expand and shrink.
- [mix-go/xutil](src/xutil) A set of tools to keep Golang sweet.

## Development Documentation

- `V1.1` https://openmix.org/mix-go/docs/1.1/
- `V1.0` https://www.kancloud.cn/onanying/mixgo1/content

## Quick Start

Provides ready-to-use scaffolding tools to quickly create projects and produce immediate output.

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

If there is an error during compilation, tidy up the dependencies.

~~~
go mod tidy
~~~

### Goland

- [How to use MixGo in IDE Goland](https://zhuanlan.zhihu.com/p/391857663)

### Examples

- [Write a CLI program](examples/cli-skeleton#readme)
  - [Write a Worker Pool queue consumer](examples/cli-skeleton#%E7%BC%96%E5%86%99%E4%B8%80%E4%B8%AA-worker-pool-%E9%98%9F%E5%88%97%E6%B6%88%E8%B4%B9)
- [Write an API service](examples/api-skeleton#readme)
- [Write a Web service](examples/web-skeleton#readme)
  - [Write a WebSocket service](examples/web-skeleton#%E7%BC%96%E5%86%99%E4%B8%80%E4%B8%AA-WebSocket-%E6%9C%8D%E5%8A%A1)
- [Write a gRPC service, client](examples/grpc-skeleton#readme)

## Technical Discussion

Zhihu: https://www.zhihu.com/people/onanying    
Official QQ Group: [284806582](https://shang.qq.com/wpa/qunwpa?idkey=b3a8618d3977cda4fed2363a666b081a31d89e3d31ab164497f53b72cf49968a), [825122875](http://shang.qq.com/wpa/qunwpa?idkey=d2908b0c7095fc7ec63a2391fa4b39a8c5cb16952f6cfc3f2ce4c9726edeaf20) Secret Password: gopher

## PHP Framework

OpenMix also has PHP ecosystem frameworks:

- https://github.com/mix-php/mix
- https://gitee.com/mix-php/mix

## License

Apache License Version 2.0, http://www.apache.org/licenses/
