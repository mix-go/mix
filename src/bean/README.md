## Mix Bean

DI、IoC 容器，参考 spring bean 设计

DI, IoC container, reference spring bean

> 该库还有 php 版本：https://github.com/mix-php/bean

## Overview

一个创建对象以及处理对象依赖关系的库，该库可以实现统一管理依赖，全局对象管理，动态配置刷新等。

## Installation

- 安装

```
go get -u github.com/mix-go/bean
```

## Usage

- `ConstructorArgs` 构造器注入

```golang
var definitions = []bean.Definition{
    {
        Name: "foo",
        Reflect: bean.NewReflect(NewHttpClient),
        ConstructorArgs: bean.ConstructorArgs{
            time.Duration(time.Second * 3),
        },
    },
}

// 必须返回指针类型
func NewHttpClient(timeout time.Duration) *http.Client {
    return &http.Client{
        Timeout: timeout,
    }
}

context := bean.NewApplicationContext(definitions)
foo := context.Get("foo").(*http.Client) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

- `Fields` 字段注入

```golang
var definitions = []bean.Definition{
    {
        Name: "foo",
        Reflect: bean.NewReflect(http.Client{}),
        Fields: bean.Fields{
            "Timeout": time.Duration(time.Second * 3),
        },
    },
}

context := bean.NewApplicationContext(definitions)
foo := context.Get("foo").(*http.Client) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

- `ConstructorArgs + Fields` 混合使用

```golang
var definitions = []bean.Definition{
    {
        Name: "foo",
        Reflect: bean.NewReflect(NewHttpClient),
        ConstructorArgs: bean.ConstructorArgs{
            time.Duration(time.Second * 3),
        },
        Fields: bean.Fields{
            "Timeout": time.Duration(time.Second * 2),
        },
    },
}

// 必须返回指针类型
func NewHttpClient(timeout time.Duration) *http.Client {
    return &http.Client{
        Timeout: timeout,
    }
}

context := bean.NewApplicationContext(definitions)
foo := context.Get("foo").(*http.Client) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

- `NewReference` 引用

引用其他依赖注入

```golang
type Foo struct {
    Client *http.Client // 引用注入的都是指针类型
}

var definitions = []bean.Definition{
    {
        Name: "foo",
        Reflect: bean.NewReflect(Foo{}),
        },
        Fields: bean.Fields{
            "Client": NewReference("bar"),
        },
    },
    {
        Name: "bar",
        Reflect: bean.NewReflect(http.Client{}),
        Fields: bean.Fields{
            "Timeout": time.Duration(time.Second * 3),
        },
    },
}

context := bean.NewApplicationContext(definitions)
foo := context.Get("foo").(*Foo) // 返回的都是指针类型
cli := foo.Client
fmt.Println(fmt.Sprintf("%+v", cli))
```

- `Scope: SINGLETON` 单例

定义组件为全局单例

```golang
var definitions = []bean.Definition{
    {
        Name: "foo",
        Scope: bean.SINGLETON, // 这里定义了单例模式
        Reflect: bean.NewReflect(http.Client{}),
        Fields: bean.Fields{
            "Timeout": time.Duration(time.Second * 3),
        },
    },
}

context := bean.NewApplicationContext(definitions)
foo := context.Get("foo").(*http.Client) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

- `InitMethod` 初始化方法

对象创建完成并且 `ConstructorArgs + Fields` 两种注入全部完成后执行该方法，用来初始化处理。

```golang
type Foo struct {
    Bar    string
}

func (c *Foo) Init() {
    c.Bar = "bar ..."
    fmt.Println("init")
}

var definitions = []bean.Definition{
    {
        Name:       "foo",
        InitMethod: "Init", // 这里定义了初始化方法
        Reflect: bean.NewReflect(Foo{}),
        Fields: bean.Fields{
            "Bar":    "bar",
        },
    },
}

context := bean.NewApplicationContext(definitions)
foo := context.Get("foo").(*Foo) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

- `Refresh` 动态刷新配置

这个通常用于通过微服务配置中心实现动态刷新微服务配置的功能。

```golang
type Foo struct {
    Bar    string
}

var definitions = []bean.Definition{
    {
        Name:       "foo",
        Reflect: bean.NewReflect(Foo{}),
        Fields: bean.Fields{
            "Bar":    "bar",
        },
    },
}

context := bean.NewApplicationContext(definitions)

// 第一次获取
foo := context.Get("foo").(*Foo) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))

// 修改配置
bd := context.GetBeanDefinition("foo")
bd.Fields["Bar"] = "bar2"
bd.Refresh()

// 第二次获取就是新的配置
foo := context.Get("foo").(*Foo) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
