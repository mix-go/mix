package bean

import (
    "github.com/stretchr/testify/assert"
    "net/http"
    "testing"
    "time"
)

type foo struct {
    Bar    string
    Client *http.Client
}

func (c *foo) Init() {
    c.Bar = "test"
}

var definitions = []BeanDefinition{
    {
        Name:    "httpclient",
        Scope:   SINGLETON,
        Reflect: NewReflect(http.Client{}),
        Fields: Fields{
            "Timeout": time.Duration(time.Second * 3),
        },
    },
    {
        Name:    "httpclient2",
        Reflect: NewReflect(NewHttpClient),
        ConstructorArgs: ConstructorArgs{
            time.Duration(time.Second * 3),
        },
        Fields: Fields{
            "Timeout": time.Duration(time.Second * 2),
        },
    },
    {
        Name:       "foo",
        InitMethod: "Init",
        Reflect:    NewReflect(foo{}),
        Fields: Fields{
            "Bar":    "bar",
            "Client": NewReference("httpclient2"),
        },
    },
}

// 只能返回指针类型方能注入成功
func NewHttpClient(timeout time.Duration) *http.Client {
    return &http.Client{
        Timeout: timeout,
    }
}

func TestApplicationContextGet(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions)
    cli := context.Get("httpclient").(*http.Client)
    _, err := cli.Get("http://www.baidu.com/")

    a.Equal(err, nil)
}

func TestApplicationContextGetReflectFunc(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions)
    cli := context.Get("httpclient2").(*http.Client)
    _, err := cli.Get("http://www.baidu.com/")

    a.Equal(err, nil)
}

func TestApplicationContextGetReflectStruct(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions)
    foo := context.Get("foo").(*foo)

    a.Equal(foo.Bar, "test")

    cli := foo.Client
    _, err := cli.Get("http://www.baidu.com/")

    a.Equal(err, nil)
}
