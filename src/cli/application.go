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
	// iApp
	iApp *Application
	// Version
	Version = "1.0.24"
)

func init() {
	iApp = New("app", "1.0.0")
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

// App
func App() *Application {
	return iApp
}

// SetName
func SetName(name string) *Application {
	return iApp.SetName(name)
}

// SetVersion
func SetVersion(version string) *Application {
	return iApp.SetVersion(version)
}

// SetDebug
func SetDebug(debug bool) *Application {
	return iApp.SetDebug(debug)
}

// Use
func Use(h ...HandlerFunc) *Application {
	return iApp.Use(h...)
}

// AddCommand
func AddCommand(cmds ...*Command) *Application {
	iApp.AddCommand(cmds...)
	return iApp
}

// Run
func Run() {
	iApp.Run()
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

	// 是否单命令
	singleton bool
	// 默认命令
	defaultCommand string
	// 命令集合
	commands []*Command
	// handlers
	handlers []HandlerFunc
}

// HandlerFunc
type HandlerFunc func(next func())

// SetName
func (t *Application) SetName(name string) *Application {
	t.Name = name
	return t
}

// SetVersion
func (t *Application) SetVersion(version string) *Application {
	t.Version = version
	return t
}

// SetDebug
func (t *Application) SetDebug(debug bool) *Application {
	t.Debug = debug
	return t
}

// Use
func (t *Application) Use(h ...HandlerFunc) *Application {
	t.handlers = append(t.handlers, h...)
	return t
}

// AddCommand
func (t *Application) AddCommand(cmds ...*Command) *Application {
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
	return t
}

// Run 执行
func (t *Application) Run() {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case *NotFoundError, *UnsupportError:
				fmt.Println(err)
				return
			default:
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
	if cmd.Run == nil && cmd.RunI == nil {
		panic(fmt.Errorf("'%s' command Run/RunI is empty", cmd.Name))
	}

	// 执行命令
	exec := func() {
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
	}
	if len(t.handlers) > 0 {
		tmp := t.handlers
		for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
			tmp[i], tmp[j] = tmp[j], tmp[i]
		}
		var next func()
		for k, f := range tmp {
			if k == 0 {
				n := exec
				c := f
				next = func() {
					c(n)
				}
			} else if len(tmp)-1 == k {
				f(next)
			} else {
				n := next
				c := f
				next = func() {
					c(n)
				}
			}
		}
	} else {
		exec()
	}
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
