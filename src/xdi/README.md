> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix XDI

DI、IoC 容器

DI, IoC container

## Overview

一个创建对象以及处理对象依赖关系的库，该库可以实现统一管理依赖，全局对象管理，动态配置刷新等。

## Installation

```
go get github.com/mix-go/xdi
```

## Quick start

通过依赖配置实例化一个单例

```go
package main

import (
    "github.com/mix-go/xdi"
)

type Foo struct {
}

func init() {
    obj := &xdi.Object{
        Name: "foo",
        New: func() (interface{}, error) {
            i := &Foo{}
            return i, nil
        },
    }
    if err := xdi.Provide(obj); err != nil {
        panic(err)
    }
}

func main() {
    var foo *Foo
    if err := xdi.Populate("foo", &foo); err != nil {
        panic(err)
    }
    // use foo
}
```

## Reference

依赖配置中引用另一个依赖配置的实例

```go
package main

import (
    "github.com/mix-go/xdi"
)

type Foo struct {
    Bar *Bar
}

type Bar struct {
}

func init() {
    objs := []*xdi.Object{
        {
            Name: "foo",
            New: func() (interface{}, error) {
                // reference bar
                var bar *Bar
                if err := xdi.Populate("bar", &bar); err != nil {
                    return nil, err
                }

                i := &Foo{
                    Bar: bar,
                }
                return i, nil
            },
        },
        {
            Name: "bar",
            New: func() (interface{}, error) {
                i := &Bar{}
                return i, nil
            },
            NewEverytime: true,
        },
    }
    if err := xdi.Provide(objs...); err != nil {
        panic(err)
    }
}

func main() {
    var foo *Foo
    if err := xdi.Populate("foo", &foo); err != nil {
        panic(err)
    }
    // use foo
}
```

## Refresh singleton

程序执行中配置信息发生变化时，`Refresh()` 可以刷新单例的实例来切换使用新的配置，通常在微服务配置中心中使用。

```go
obj, err := xdi.Container().Object("foo")
if err != nil {
    panic(err)
}
if err := obj.Refresh(); err != nil {
    panic(err)
}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
