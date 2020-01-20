package bean

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "reflect"
    "testing"
    "time"
)

type Foo struct {
    Bar    string
    Client *http.Client
}

func (c *Foo) Init() {
    c.Bar = "bar init"
    fmt.Println("init")
}

var Definitions = []Definition{
    {
        Name:  "httpclient",
        Scope: SINGLETON,
        Reflect: func() reflect.Value {
            return reflect.New(reflect.TypeOf(http.Client{}))
        },
        Fields: Fields{
            "Timeout": time.Duration(time.Second * 3),
        },
    },
    {
        Name: "httpclient2",
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
    {
        Name:       "foo",
        InitMethod: "Init",
        Reflect: func() reflect.Value {
            return reflect.New(reflect.TypeOf(Foo{}))
        },
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

func TestApplicationContext_Get(t *testing.T) {
    context := NewApplicationContext(Definitions)
    cli := context.Get("httpclient").(*http.Client)
    fmt.Println(fmt.Sprintf("%#v", cli))
    resp, _ := cli.Get("http://www.baidu.com/")
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(len(string(body)))
}

func TestApplicationContext_Get2(t *testing.T) {
    context := NewApplicationContext(Definitions)
    cli := context.Get("httpclient2").(*http.Client)
    fmt.Println(fmt.Sprintf("%#v", cli))
    resp, _ := cli.Get("http://www.baidu.com/")
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(len(string(body)))
}

func TestApplicationContext_Get3(t *testing.T) {
    context := NewApplicationContext(Definitions)
    foo := context.Get("foo").(*Foo)
    fmt.Println(fmt.Sprintf("%#v", foo))
    cli := foo.Client
    fmt.Println(fmt.Sprintf("%#v", cli))
    resp, _ := cli.Get("http://www.baidu.com/")
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(len(string(body)))
}
