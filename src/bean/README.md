## Mix Bean

DI、ICO 容器，参考 spring bean 设计

DI, ICO container, reference spring bean

> 该库还有 php 版本：https://github.com/mix-php/bean

## Usage

- 安装

```
go get -u github.com/mix-go/bean
```

- 构造器注入

```golang
var Definitions = []Definition{
    {
        Name: "foo",
        Reflect: func() reflect.Value {
            return reflect.ValueOf(NewHttpClient)
        },
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

- 字段注入

```golang
var Definitions = []Definition{
    {
        Name: "foo",
        Reflect: func() reflect.Value {
            return reflect.New(reflect.TypeOf(http.Client{}))
        },
        Fields: Fields{
            "Timeout": time.Duration(time.Second * 3),
        },
    },
}

context := NewApplicationContext(Definitions)
foo := context.Get("foo").(*http.Client) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

- 混合使用

```golang
var Definitions = []Definition{
    {
        Name: "foo",
        Reflect: func() reflect.Value {
            return reflect.ValueOf(NewHttpClient)
        },
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

- 引用

```golang
type Foo struct {
    Client *http.Client // 引用注入的都是指针类型
}

var Definitions = []Definition{
    {
        Name: "foo",
        Reflect: func() reflect.Value {
            return reflect.New(reflect.TypeOf(Foo{}))
        },
        Fields: Fields{
            "Client": Reference{Name: "bar"},
        },
    },
    {
        Name: "bar",
        Reflect: func() reflect.Value {
            return reflect.New(reflect.TypeOf(http.Client{}))
        },
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

- 单例

```golang
var Definitions = []Definition{
    {
        Name: "foo",
        Scope: SINGLETON,
        Reflect: func() reflect.Value {
            return reflect.New(reflect.TypeOf(http.Client{}))
        },
        Fields: Fields{
            "Timeout": time.Duration(time.Second * 3),
        },
    },
}

context := NewApplicationContext(Definitions)
foo := context.Get("foo").(*http.Client) // 返回的都是指针类型
fmt.Println(fmt.Sprintf("%+v", foo))
```

- 初始化方法

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
        Reflect: func() reflect.Value {
            return reflect.New(reflect.TypeOf(Foo{}))
        },
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
