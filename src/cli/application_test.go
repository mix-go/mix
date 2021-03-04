package cli

import (
	"fmt"
	"github.com/mix-go/cli/argv"
	"github.com/mix-go/cli/flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var ok = false

func TestCommandRun(t *testing.T) {
	a := assert.New(t)

	os.Args = []string{os.Args[0], "foo"}
	argv.Parse()
	flag.Parse()

	cmd := &Command{
		Name:  "foo",
		Usage: "bar",
		Run: func() {
			ok = true
		},
	}
	opt := &Option{
		Names: []string{"a", "bc"},
		Usage: "foo",
	}
	cmd.AddOption(opt)
	app := NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()

	a.NotEqual(app.BasePath, nil)
	a.True(ok)
	ok = false
}

func TestSingletonCommandRun(t *testing.T) {
	a := assert.New(t)

	os.Args = []string{os.Args[0], "-a"}
	argv.Parse()
	flag.Parse()

	cmd := &Command{
		Name:  "foo",
		Usage: "bar",
		Run: func() {
			ok = true
		},
		Singleton: true,
	}
	app := NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()

	a.NotEqual(app.BasePath, nil)
	a.True(ok)
	ok = false
}

func TestDefaultCommandRun(t *testing.T) {
	a := assert.New(t)

	// 多命令
	os.Args = []string{os.Args[0]}
	argv.Parse()
	flag.Parse()

	cmd := &Command{
		Name:  "foo",
		Usage: "bar",
		Run: func() {
			ok = true
		},
		Default: true,
	}
	app := NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()

	a.NotEqual(app.BasePath, nil)
	a.True(ok)
	ok = false

	// 单命令
	os.Args = []string{os.Args[0]}
	argv.Parse()
	flag.Parse()

	cmd = &Command{
		Name:  "foo",
		Usage: "bar",
		Run: func() {
			ok = true
		},
		Singleton: true,
		Default:   true,
	}
	app = NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()

	a.NotEqual(app.BasePath, nil)
	a.True(ok)
	ok = false
}

func TestCommandNotFound(t *testing.T) {
	os.Args = []string{os.Args[0], "bar"}
	argv.Parse()
	flag.Parse()

	cmd := &Command{
		Name:  "foo",
		Usage: "bar",
		Run: func() {
			ok = true
		},
	}
	app := NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()

	fmt.Println("-----------------------")

	// 默认 + 找不到
	os.Args = []string{os.Args[0], "中文foo"}
	argv.Parse()
	flag.Parse()

	cmd = &Command{
		Name:  "foo",
		Usage: "bar",
		Run: func() {
			ok = true
		},
		Default: true,
	}
	app = NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()
}

func TestCommandPrint(t *testing.T) {
	os.Args = []string{os.Args[0]}
	fmt.Println(os.Args)
	argv.Parse()
	flag.Parse()
	cmd := &Command{
		Name:  "foo",
		Usage: "bar",
		Run: func() {
		},
	}
	app := NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()

	fmt.Println("-----------------------")

	os.Args = []string{os.Args[0], "-h"}
	fmt.Println(os.Args)
	argv.Parse()
	flag.Parse()
	app = NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()

	fmt.Println("-----------------------")

	os.Args = []string{os.Args[0], "-v"}
	fmt.Println(os.Args)
	argv.Parse()
	flag.Parse()
	app = NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()

	fmt.Println("-----------------------")

	os.Args = []string{os.Args[0], "foo", "--help"}
	fmt.Println(os.Args)
	argv.Parse()
	flag.Parse()
	app = NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()

	fmt.Println("-----------------------")

	os.Args = []string{os.Args[0]}
	fmt.Println(os.Args)
	argv.Parse()
	flag.Parse()
	cmd = &Command{
		Name:  "foo",
		Usage: "bar",
		Run: func() {
		},
		Singleton: true,
	}
	app = NewApplication("test", "1.0.0")
	app.AddCommand(cmd)
	app.Run()
}
