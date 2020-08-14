package console

import (
    "fmt"
    "github.com/mix-go/bean"
    "github.com/mix-go/console/argv"
    "github.com/mix-go/console/flag"
    "github.com/mix-go/event"
    "github.com/mix-go/logrus"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

var (
    def1 = ApplicationDefinition{
        AppName:    "app-test",
        AppVersion: "1.0.0-test",
        AppDebug:   true,
        Beans: []bean.BeanDefinition{
            {
                Name:    "eventDispatcher",
                Reflect: bean.NewReflect(event.NewDispatcher),
                Scope:   bean.SINGLETON,
            },
            {
                Name:            "error",
                Reflect:         bean.NewReflect(NewError),
                Scope:           bean.SINGLETON,
                ConstructorArgs: bean.ConstructorArgs{bean.NewReference("logger")},
            },
            {
                Name:    "logger",
                Reflect: bean.NewReflect(logrus.NewLogger),
                Scope:   bean.SINGLETON,
            },
        },
        Commands: []CommandDefinition{
            {
                Name:  "foo",
                Usage: "bar",
                Options: []OptionDefinition{
                    {
                        Names: []string{"a", "bc"},
                        Usage: "foo",
                    },
                },
                Reflect:   bean.NewReflect(Foo{}),
                Singleton: false,
            },
        },
    }

    def2 = ApplicationDefinition{
        AppName:    "app-test",
        AppVersion: "1.0.0-test",
        AppDebug:   true,
        Beans: []bean.BeanDefinition{
            {
                Name:    "eventDispatcher",
                Reflect: bean.NewReflect(event.NewDispatcher),
                Scope:   bean.SINGLETON,
            },
            {
                Name:            "error",
                Reflect:         bean.NewReflect(NewError),
                Scope:           bean.SINGLETON,
                ConstructorArgs: bean.ConstructorArgs{bean.NewReference("logger")},
            },
            {
                Name:    "logger",
                Reflect: bean.NewReflect(logrus.NewLogger),
                Scope:   bean.SINGLETON,
            },
        },
        Commands: []CommandDefinition{
            {
                Name:  "foo",
                Usage: "bar",
                Options: []OptionDefinition{
                    {
                        Names: []string{"a", "bc"},
                        Usage: "foo",
                    },
                },
                Reflect:   bean.NewReflect(Foo{}),
                Singleton: true,
            },
        },
    }

    run = false
)

type Foo struct {
    Bar string
}

func (c *Foo) Main() {
    run = true
}

func TestCommandRun(t *testing.T) {
    a := assert.New(t)

    os.Args = []string{os.Args[0], "foo"}
    argv.Parse()
    flag.Parse()

    app := NewApplication(def1, "eventDispatcher", "error");
    app.Run()

    a.NotEqual(app.BasePath, nil)
    a.True(run)

    run = false
}

func TestSingletonCommandRun(t *testing.T) {
    a := assert.New(t)

    os.Args = []string{os.Args[0], "-a"}
    argv.Parse()
    flag.Parse()

    app := NewApplication(def2, "eventDispatcher", "error");
    app.Run()

    a.NotEqual(app.BasePath, nil)
    a.True(run)

    run = false
}

func TestCommandNotFound(t *testing.T) {
    a := assert.New(t)

    os.Args = []string{os.Args[0], "bar"}
    argv.Parse()
    flag.Parse()
    app := NewApplication(def1, "eventDispatcher", "error")
    app.Run()

    a.Contains(LastError.(error).Error(), "'bar' is not command, see '")
}

func TestCommandPrint(t *testing.T) {
    var app *Application

    os.Args = []string{os.Args[0]}
    fmt.Println(os.Args)
    argv.Parse()
    flag.Parse()
    app = NewApplication(def1, "eventDispatcher", "error")
    app.Run()

    fmt.Println("-----------------------")

    os.Args = []string{os.Args[0], "-h"}
    fmt.Println(os.Args)
    argv.Parse()
    flag.Parse()
    app = NewApplication(def1, "eventDispatcher", "error")
    app.Run()

    fmt.Println("-----------------------")

    os.Args = []string{os.Args[0], "-v"}
    fmt.Println(os.Args)
    argv.Parse()
    flag.Parse()
    app = NewApplication(def1, "eventDispatcher", "error")
    app.Run()

    fmt.Println("-----------------------")

    os.Args = []string{os.Args[0], "foo", "--help"}
    fmt.Println(os.Args)
    argv.Parse()
    flag.Parse()
    app = NewApplication(def1, "eventDispatcher", "error")
    app.Run()

    fmt.Println("-----------------------")

    os.Args = []string{os.Args[0]}
    fmt.Println(os.Args)
    argv.Parse()
    flag.Parse()
    app = NewApplication(def2, "eventDispatcher", "error")
    app.Run()
}
