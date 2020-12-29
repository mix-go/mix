package console

import (
	"errors"
	"fmt"
	"github.com/mix-go/bean"
	"github.com/mix-go/console/argv"
	"github.com/mix-go/console/flag"
	"github.com/mix-go/event"
	"os"
	"strings"
)

var (
	// App 全局APP
	App *Application
	// Version 版本号
	Version = "1.0.24"
	// LastError 最后的错误
	LastError interface{}
)

// NewApplication 创建应用
func NewApplication(definition ApplicationDefinition, dispatcherName, errorName string) *Application {
	App = &Application{
		ApplicationDefinition: definition,
		DispatcherName:        dispatcherName,
		ErrorName:             errorName,
	}
	App.Init()
	return App
}

// ApplicationDefinition 应用定义
type ApplicationDefinition struct {
	// 应用名称
	Name string
	// 应用版本
	Version string
	// 应用调试
	Debug bool
	// 依赖配置
	Beans []bean.BeanDefinition
	// 命令集合
	Commands []CommandDefinition
}

// Application 应用
type Application struct {
	// App 定义
	ApplicationDefinition
	// Event Dispatcher
	DispatcherName string
	Dispatcher     event.Dispatcher
	// Error
	ErrorName string
	Error     Error
	// 基础路径
	BasePath string
	// 应用上下文
	Context *bean.ApplicationContext
	// 是否单命令
	Singleton bool
	// 默认命令
	DefaultCommand string
}

// CommandDefinition 命令定义
type CommandDefinition struct {
	// 命令名称
	Name string
	// 使用描述
	Usage string
	// 选项
	Options []OptionDefinition
	// 命令
	Command Command
	// 是否单命令
	Singleton bool
	// 是否为默认命令
	Default bool
}

// Command 命令接口
type Command interface {
	Main()
}

// OptionDefinition 命令选项
type OptionDefinition struct {
	Names []string
	Usage string
}

// Init 初始化
func (t *Application) Init() {
	t.Context = bean.NewApplicationContext(t.Beans)

	t.Dispatcher = t.Context.Get(t.DispatcherName).(event.Dispatcher)
	t.Error = t.Context.Get(t.ErrorName).(Error)

	t.BasePath = argv.Program().Dir

	for _, c := range t.Commands {
		if c.Singleton {
			t.Singleton = true
		}
		if c.Default {
			t.DefaultCommand = c.Name
		}
	}
	if t.Singleton {
		argv.Parse(true)
		flag.Parse()
	}
}

// Get 快速获取实例
func (t *Application) Get(name string) interface{} {
	return t.Context.Get(name)
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

			t.Error.Handle(err)
		}
	}()

	if len(t.Commands) == 0 {
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
			if t.DefaultCommand != "" && len(os.Args) == 1 {
				os.Args = append(os.Args, t.DefaultCommand)
				argv.Parse()
				flag.Parse()
				t.Run()
			} else {
				t.globalHelp()
			}
			return
		}

		if t.Singleton {
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
	var d *CommandDefinition
	command := argv.Command()
	if t.Singleton {
		// 单命令
		for _, c := range t.Commands {
			if c.Singleton {
				d = &c
				break
			}
		}
		if d == nil {
			panic(errors.New("singleton command not found"))
		}
	} else {
		for _, c := range t.Commands {
			if c.Name == command {
				d = &c
				break
			}
		}
	}
	if d == nil {
		panic(NewNotFoundError(fmt.Errorf("'%s' is not command, see '%s --help'.", command, argv.Program().Path)))
	}

	// 获取命令
	c := d.Command

	// 触发执行命令前置事件
	e := &CommandBeforeExecuteEvent{
		Command: c,
	}
	t.Dispatcher.Dispatch(e)

	// 执行命令
	c.Main()
}

// 命令行选项效验
func (t *Application) validateOptions() {
	var options []OptionDefinition
	if !t.Singleton {
		for _, v := range t.Commands {
			if v.Name == argv.Command() {
				options = v.Options
				break
			}
		}
	} else {
		for _, v := range t.Commands {
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
	if !t.Singleton {
		fg = " [OPTIONS] COMMAND"
	}
	fmt.Println(fmt.Sprintf("Usage: %s%s [opt...]", program, fg))
	if !t.Singleton {
		t.printCommands()
	} else {
		t.printCommandOptions()
	}
	t.printGlobalOptions()
	fmt.Println("")
	fg = ""
	if !t.Singleton {
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
	for _, v := range t.Commands {
		command := v.Name
		usage := v.Usage
		fmt.Println(fmt.Sprintf("  %s\t%s", command, usage))
	}
}

// 打印命令选项
func (t *Application) printCommandOptions() {
	var options []OptionDefinition
	if !t.Singleton {
		for _, v := range t.Commands {
			if v.Name == argv.Command() {
				options = v.Options
				break
			}
		}
	} else {
		for _, v := range t.Commands {
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
