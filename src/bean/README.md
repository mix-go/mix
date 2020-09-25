## Mix Bean

DI、IoC 容器，参考 spring bean 设计

DI, IoC container, reference spring bean

> 该库还有 php 版本：https://github.com/mix-php/bean

## Installation

- 安装

```
go get -u github.com/mix-go/bean
```

## Usage

- `ConstructorArgs` 构造器注入

```golang
var Definitions = []Definition{
    {
        Name: "foo",
        Reflect: NewReflect(NewHttpClient),
        ConstructorArgs: ConstructorArgs{
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

context := NewApplicationContext(Definitions)
foo := context.Get("foo").(*http.Client) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

- `Fields` 字段注入

```golang
var Definitions = []Definition{
    {
        Name: "foo",
        Reflect: NewReflect(http.Client{}),
        Fields: Fields{
            "Timeout": time.Duration(time.Second * 3),
        },
    },
}

context := NewApplicationContext(Definitions)
foo := context.Get("foo").(*http.Client) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

- `ConstructorArgs + Fields` 混合使用

```golang
var Definitions = []Definition{
    {
        Name: "foo",
        Reflect: NewReflect(NewHttpClient),
        ConstructorArgs: ConstructorArgs{
            time.Duration(time.Second * 3),
        },
        Fields: Fields{
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

context := NewApplicationContext(Definitions)
foo := context.Get("foo").(*http.Client) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

- `NewReference` 引用

引用其他依赖注入

```golang
type Foo struct {
    Client *http.Client // 引用注入的都是指针类型
}

var Definitions = []Definition{
    {
        Name: "foo",
        Reflect: NewReflect(Foo{}),
        },
        Fields: Fields{
            "Client": NewReference("bar"),
        },
    },
    {
        Name: "bar",
        Reflect: NewReflect(http.Client{}),
        Fields: Fields{
            "Timeout": time.Duration(time.Second * 3),
        },
    },
}

context := NewApplicationContext(Definitions)
foo := context.Get("foo").(*Foo) // 返回的都是指针类型
cli := foo.Client
fmt.Println(fmt.Sprintf("%+v", cli))
```

- `Scope: SINGLETON` 单例

定义组件为全局单例

```golang
var Definitions = []Definition{
    {
        Name: "foo",
        Scope: SINGLETON,
        Reflect: NewReflect(http.Client{}),
        Fields: Fields{
            "Timeout": time.Duration(time.Second * 3),
        },
    },
}

context := NewApplicationContext(Definitions)
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

var Definitions = []Definition{
    {
        Name:       "foo",
        InitMethod: "Init",
        Reflect: NewReflect(Foo{}),
        Fields: Fields{
            "Bar":    "bar",
        },
    },
}

context := NewApplicationContext(Definitions)
foo := context.Get("foo").(*Foo) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
