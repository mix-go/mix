package console

import (
    "github.com/mix-go/bean"
    "github.com/mix-go/console/cli"
    "github.com/mix-go/console/flag"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

var (
    def1 = ApplicationDefinition{
        AppName:    "app-test",
        AppVersion: "1.0.0-test",
        AppDebug:   true,
        Beans:      nil,
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
    def2 = ApplicationDefinition{
        AppName:    "app-test",
        AppVersion: "1.0.0-test",
        AppDebug:   true,
        Beans:      nil,
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
                Singleton: true,
            },
        },
    }
    Run = false
)

type Foo struct {
    Bar string
}

func (c *Foo) Main() {
    Run = true
}

func TestCommandRun(t *testing.T) {
    a := assert.New(t)

    os.Args = []string{os.Args[0], "foo"}
    cli.Parse()
    flag.Parse()

    app := NewApplication(def1);
    app.Run()

    a.NotEqual(app.BasePath, nil)
    a.True(Run)

    Run = false
}

func TestSingletonCommandRun(t *testing.T) {
    a := assert.New(t)

    os.Args = []string{os.Args[0]}
    cli.Parse()
    flag.Parse()

    app := NewApplication(def2);
    app.Run()

    a.NotEqual(app.BasePath, nil)
    a.True(Run)

    Run = false
}
