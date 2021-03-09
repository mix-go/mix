package cli

import (
	"errors"
	"fmt"
	"github.com/mix-go/cli/argv"
	"github.com/mix-go/cli/flag"
	"os"
	"strings"
)

var (
	// App
	App *Application
	// Version
	Version = "1.0.24"
)

func init() {
	App = New("app", "1.0.0")
}

// New
func New(name, version string) *Application {
	app := &Application{
		Name:    name,
		Version: version,
	}
	app.BasePath = argv.Program().Dir
	return app
}

// SetName
func SetName(name string) *Application {
	App.Name = name
	return App
}

// SetVersion
func SetVersion(version string) *Application {
	App.Version = version
	return App
}

// SetDebug
func SetDebug(debug bool) *Application {
	App.Debug = debug
	return App
}

// AddCommand
func AddCommand(cmds ...*Command) *Application {
	App.AddCommand(cmds...)
	return App
}

// Run
func Run() {
	App.Run()
}

// Application
type Application struct {
	// 应用名称
	Name string
	// 应用版本
	Version string
	// 应用调试
	Debug bool
	// 基础路径
	BasePath string
	// ErrorHandle
	ErrorHandle func(err interface{})

	// 是否单命令
	singleton bool
	// 默认命令
	defaultCommand string
	// 命令集合
	commands []*Command
}

// AddCommand
func (t *Application) AddCommand(cmds ...*Command) {
	t.commands = append(t.commands, cmds...)
	// init
	for _, c := range t.commands {
		if c.Singleton {
			t.singleton = true
		}
		if c.Default {
			t.defaultCommand = c.Name
		}
	}
	if t.singleton {
		argv.Parse(true)
		flag.Parse()
	}
}

// Run 执行
func (t *Application) Run() {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case *NotFoundError, *UnsupportError:
				fmt.Println(err)
				return
			}
			if t.ErrorHandle != nil {
				t.ErrorHandle(err)
			} else {
				panic(err)
			}
		}
	}()

	if len(t.commands) == 0 {
		panic(errors.New("command cannot be empty"))
	}

	command := argv.Command()
	if command == "" {
		if flag.Match("h", "help").Bool(false) {
			t.globalHelp()
			return
		}
		if flag.Match("v", "version").Bool(false) {
			t.version()
			return
		}

		options := flag.Options().Map()
		if len(options) == 0 {
			if t.defaultCommand != "" && len(os.Args) == 1 {
				os.Args = append(os.Args, t.defaultCommand)
				argv.Parse()
				flag.Parse()
				t.Run()
			} else {
				t.globalHelp()
			}
			return
		}

		if t.singleton {
			t.call()
			return
		}

		f := ""
		for k := range options {
			f = k
			break
		}
		p := argv.Program().Path
		panic(NewNotFoundError(fmt.Errorf("flag provided but not defined: '%s', see '%s --help'.", f, p)))
	} else if flag.Match("help").Bool(false) {
		t.commandHelp()
		return
	}
	t.call()
}

// 执行命令
func (t *Application) call() {
	// 命令行选项效验
	t.validateOptions()

	// 提取命令
	var cmd *Command
	command := argv.Command()
	if t.singleton {
		// 单命令
		for _, c := range t.commands {
			if c.Singleton {
				cmd = c
				break
			}
		}
		if cmd == nil {
			panic(errors.New("singleton command not found"))
		}
	} else {
		for _, c := range t.commands {
			if c.Name == command {
				cmd = c
				break
			}
		}
	}
	if cmd == nil {
		panic(NewNotFoundError(fmt.Errorf("'%s' is not command, see '%s --help'.", command, argv.Program().Path)))
	}

	// 执行命令
	r := cmd.Run
	if r != nil {
		r()
		return
	}
	ri := cmd.RunI
	if ri != nil {
		ri.Main()
		return
	}
	panic(fmt.Errorf("'%s' command Run/RunI is empty", cmd.Name))
}

// 命令行选项效验
func (t *Application) validateOptions() {
	var options []*Option
	if !t.singleton {
		for _, v := range t.commands {
			if v.Name == argv.Command() {
				options = v.Options
				break
			}
		}
	} else {
		for _, v := range t.commands {
			if v.Singleton {
				options = v.Options
				break
			}
		}
	}
	if len(options) == 0 {
		return
	}

	var flags []string
	for _, o := range options {
		for _, v := range o.Names {
			if len(v) == 1 {
				flags = append(flags, fmt.Sprintf("-%s", v))
			} else {
				flags = append(flags, fmt.Sprintf("--%s", v))
			}
		}
	}
	inArray := func(value string, values []string) bool {
		for _, v := range values {
			if v == value {
				return true
			}
		}
		return false
	}
	for f := range flag.Options().Map() {
		if !inArray(f, flags) {
			p := argv.Program().Path
			c := argv.Command()
			if c != "" {
				c = fmt.Sprintf(" %s", c)
			}
			panic(NewNotFoundError(fmt.Errorf("flag provided but not defined: '%s', see '%s%s --help'.", f, p, c)))
		}
	}
}

// 全局帮助
func (t *Application) globalHelp() {
	program := argv.Program().Path
	fg := ""
	if !t.singleton {
		fg = " [OPTIONS] COMMAND"
	}
	fmt.Println(fmt.Sprintf("Usage: %s%s [opt...]", program, fg))
	if !t.singleton {
		t.printCommands()
	} else {
		t.printCommandOptions()
	}
	t.printGlobalOptions()
	fmt.Println("")
	fg = ""
	if !t.singleton {
		fg = " COMMAND"
	}
	fmt.Println("")
	fmt.Println(fmt.Sprintf("Run '%s%s --help' for more information on a command.", program, fg))
	fmt.Println("")
	fmt.Println("Developed with Mix Go framework. (openmix.org/mix-go)")
}

// 命令帮助
func (t *Application) commandHelp() {
	program := argv.Program().Path
	command := argv.Command()
	fmt.Println(fmt.Sprintf("Usage: %s %s [opt...]", program, command))
	t.printCommandOptions()
	fmt.Println("")
	fmt.Println("Developed with Mix Go framework. (openmix.org/mix-go)")
}

// 打印全局选项
func (t *Application) printGlobalOptions() {
	tabs := "\t"
	fmt.Println("")
	fmt.Println("Global Options:")
	fmt.Println(fmt.Sprintf("  -h, --help%sPrint usage", tabs))
	fmt.Println(fmt.Sprintf("  -v, --version%sPrint version information", tabs))
}

// 打印命令
func (t *Application) printCommands() {
	fmt.Println("")
	fmt.Println("Commands:")
	for _, v := range t.commands {
		command := v.Name
		usage := v.Usage
		fmt.Println(fmt.Sprintf("  %s\t%s", command, usage))
	}
}

// 打印命令选项
func (t *Application) printCommandOptions() {
	var options []*Option
	if !t.singleton {
		for _, v := range t.commands {
			if v.Name == argv.Command() {
				options = v.Options
				break
			}
		}
	} else {
		for _, v := range t.commands {
			if v.Singleton {
				options = v.Options
				break
			}
		}
	}
	if len(options) == 0 {
		return
	}

	fmt.Println("")
	fmt.Println("Command Options:")
	for _, o := range options {
		var flags []string
		for _, v := range o.Names {
			if len(v) == 1 {
				flags = append(flags, fmt.Sprintf("-%s", v))
			} else {
				flags = append(flags, fmt.Sprintf("--%s", v))
			}
		}
		fg := strings.Join(flags, ", ")
		usage := o.Usage
		fmt.Println(fmt.Sprintf("  %s\t%s", fg, usage))
	}
}

// 版本号
func (t *Application) version() {
	appName := t.Name
	appVersion := t.Version
	frameworkVersion := Version
	fmt.Println(fmt.Sprintf("%s %s, framework %s", appName, appVersion, frameworkVersion))
}
