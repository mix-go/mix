package console

import (
    "fmt"
    "github.com/mix-go/bean"
    "github.com/mix-go/console/cli"
    "github.com/mix-go/console/flag"
    "github.com/stretchr/testify/assert"
    "net/http"
    "os"
    "testing"
    "time"
)

var (
    def = ApplicationDefinition{
        AppName:    "app-test",
        AppVersion: "1.0.0-test",
        AppDebug:   true,
        Beans: []bean.BeanDefinition{
            {
                Name:    "httpclient",
                Scope:   bean.SINGLETON,
                Reflect: bean.NewReflect(http.Client{}),
                Fields: bean.Fields{
                    "Timeout": time.Duration(time.Second * 3),
                },
            },
        },
        Commands: []CommandDefinition{
            {
                Name:  "foo",
                Usage: "bar",
                Options: []OptionDefinition{
                    {
                        Names: []string{"", ""},
                        Usage: "",
                    },
                },
                Reflect:   bean.NewReflect(Foo{}),
                Singleton: false,
            },
        },
    }
)

type Foo struct {
    Bar string
}

func (c *Foo) Main() {
    fmt.Printf("%#v\n", flag.Options)
}

func TestApplication(t *testing.T) {
    a := assert.New(t)

    os.Args = []string{os.Args[0], "foo", "-a", "=", "a1", "-b", "--ab", "=", "ab1"}
    cli.Refresh()
    flag.Refresh()

    app := NewApplication(def);
    app.Run()

    a.NotEqual(app.BasePath, nil)
}
