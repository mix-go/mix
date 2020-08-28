package bean

import (
    "fmt"
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
    },
    {
        Name:    "httpclient3",
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

func TestGetReflectStruct(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions)
    cli := context.Get("httpclient").(*http.Client)
    _, err := cli.Get("http://www.baidu.com/")

    a.Equal(err, nil)
}

func TestGetReflectFunc(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions)
    cli := context.Get("httpclient3").(*http.Client)
    _, err := cli.Get("http://www.baidu.com/")

    a.Equal(err, nil)
}

func TestGetReference(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions)
    foo := context.Get("foo").(*foo)

    a.Equal(foo.Bar, "test")

    cli := foo.Client
    _, err := cli.Get("http://www.baidu.com/")

    a.Equal(err, nil)
}

func TestCoverConstructorArgs(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions)
    cli := context.Get("httpclient2").(*http.Client)
    cli1 := context.GetBean("httpclient2", nil, nil).(*http.Client)

    a.Equal(fmt.Sprintf("%v", cli), fmt.Sprintf("%v", cli1))

    cli2 := context.GetBean("httpclient2", nil, ConstructorArgs{time.Duration(time.Second * 4)}).(*http.Client)

    a.NotEqual(fmt.Sprintf("%v", cli), fmt.Sprintf("%v", cli2))
}

func TestCoverFields(t *testing.T) {
    a := assert.New(t)

    context := NewApplicationContext(definitions)
    cli := context.Get("httpclient").(*http.Client)
    cli1 := context.GetBean("httpclient", nil, nil).(*http.Client)

    a.Equal(fmt.Sprintf("%v", cli), fmt.Sprintf("%v", cli1))

    cli2 := context.GetBean("httpclient", Fields{"Timeout": time.Duration(time.Second * 4)}, nil).(*http.Client)

    a.NotEqual(fmt.Sprintf("%v", cli), fmt.Sprintf("%v", cli2))
}
